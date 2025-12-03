package controllers

import (
	"strconv"

	"controller_api_demo/services"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/web"
)

// UserController handles user-related HTTP requests.
// This follows the ASP.NET Core Controller pattern.
type UserController struct {
	// Services will be injected when the controller is created
	userService services.UserService
}

// NewUserController creates a new UserController.
// Dependencies are passed as parameters (constructor injection pattern).
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// MapRoutes registers all routes for this controller.
// This is similar to ASP.NET Core's attribute routing.
// Implements web.IController interface.
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
	// Create a route group for /api/users
	users := app.MapGroup("/api/users")
	users.WithTags("Users")

	// GET /api/users
	users.MapGet("", ctrl.ListUsers).
		WithSummary("List all users").
		WithDescription("Returns a list of all users in the system")

	// GET /api/users/{id}
	users.MapGet("/{id}", ctrl.GetUser).
		WithSummary("Get user by ID").
		WithDescription("Returns a single user by their ID")

	// POST /api/users
	users.MapPost("", ctrl.CreateUser).
		WithSummary("Create a new user").
		WithDescription("Creates a new user in the system")

	// PUT /api/users/{id}
	users.MapPut("/{id}", ctrl.UpdateUser).
		WithSummary("Update user").
		WithDescription("Updates an existing user")

	// DELETE /api/users/{id}
	users.MapDelete("/{id}", ctrl.DeleteUser).
		WithSummary("Delete user").
		WithDescription("Deletes a user from the system")
}

// ListUsers handles GET /api/users
func (ctrl *UserController) ListUsers(c *gin.Context) {
	users, err := ctrl.userService.ListUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
}

// GetUser handles GET /api/users/{id}
func (ctrl *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctrl.userService.GetUser(idInt)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

// CreateUser handles POST /api/users
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user services.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userService.CreateUser(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, user)
}

// UpdateUser handles PUT /api/users/{id}
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var user services.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user.ID = idInt
	if err := ctrl.userService.UpdateUser(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

// DeleteUser handles DELETE /api/users/{id}
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := ctrl.userService.DeleteUser(idInt); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, nil)
}
