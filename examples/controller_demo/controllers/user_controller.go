package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// User represents a user entity.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService provides user operations.
type UserService interface {
	GetUser(id int) (*User, error)
	ListUsers() ([]*User, error)
	CreateUser(user *User) error
	UpdateUser(id int, user *User) error
	DeleteUser(id int) error
}

// userService is the implementation.
type userService struct {
	users map[int]*User
}

// NewUserService creates a new UserService.
func NewUserService() UserService {
	return &userService{
		users: map[int]*User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
	}
}

func (s *userService) GetUser(id int) (*User, error) {
	user, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (s *userService) ListUsers() ([]*User, error) {
	users := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users, nil
}

func (s *userService) CreateUser(user *User) error {
	s.users[user.ID] = user
	return nil
}

func (s *userService) UpdateUser(id int, user *User) error {
	s.users[id] = user
	return nil
}

func (s *userService) DeleteUser(id int) error {
	delete(s.users, id)
	return nil
}

// ============================================
// UserController - Controller Pattern
// ============================================

// UserController handles user-related HTTP requests.
// This follows the .NET Controller pattern.
type UserController struct {
	app         *web.WebApplication
	userService UserService
}

// NewUserController creates a new UserController.
func NewUserController(app *web.WebApplication) *UserController {
	// Resolve service from DI container
	userService := di.GetRequiredService[UserService](app.Services)

	return &UserController{
		app:         app,
		userService: userService,
	}
}

// MapRoutes registers all routes for this controller.
// This is similar to .NET's attribute routing.
// Implements web.IController interface.
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
	// Create a route group for /api/users
	users := app.MapGroup("/api/users")
	users.WithTags("Users")

	// GET /api/users
	users.MapGet("", ctrl.GetAll).
		WithSummary("Get all users").
		WithDescription("Returns a list of all users in the system")

	// GET /api/users/{id}
	users.MapGet("/{id}", ctrl.GetByID).
		WithSummary("Get user by ID").
		WithDescription("Returns a single user by their ID")

	// POST /api/users
	users.MapPost("", ctrl.Create).
		WithSummary("Create a new user").
		WithDescription("Creates a new user in the system")

	// PUT /api/users/{id}
	users.MapPut("/{id}", ctrl.Update).
		WithSummary("Update a user").
		WithDescription("Updates an existing user")

	// DELETE /api/users/{id}
	users.MapDelete("/{id}", ctrl.Delete).
		WithSummary("Delete a user").
		WithDescription("Deletes a user from the system")
}

// ============================================
// Action Methods (similar to .NET Controller actions)
// ============================================

// GetAll handles GET /api/users
func (ctrl *UserController) GetAll(c *gin.Context) {
	users, err := ctrl.userService.ListUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
}

// GetByID handles GET /api/users/{id}
func (ctrl *UserController) GetByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := ctrl.userService.GetUser(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

// Create handles POST /api/users
func (ctrl *UserController) Create(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userService.CreateUser(&user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

// Update handles PUT /api/users/{id}
func (ctrl *UserController) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userService.UpdateUser(idInt, &user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

// Delete handles DELETE /api/users/{id}
func (ctrl *UserController) Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ctrl.userService.DeleteUser(idInt); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, nil)
}

