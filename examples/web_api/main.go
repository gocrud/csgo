package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/openapi"
	"github.com/gocrud/csgo/swagger"
	"github.com/gocrud/csgo/web"
)

// User represents a user model.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserRequest represents a request to create a user.
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UserService provides user-related operations.
type UserService struct {
	users  map[int]*User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users: map[int]*User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (s *UserService) GetUser(id int) *User {
	return s.users[id]
}

func (s *UserService) ListUsers() []*User {
	users := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

func (s *UserService) CreateUser(name, email string) *User {
	user := &User{ID: s.nextID, Name: name, Email: email}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserService) DeleteUser(id int) bool {
	if _, ok := s.users[id]; ok {
		delete(s.users, id)
		return true
	}
	return false
}

func main() {
	// Create web application builder (corresponds to .NET WebApplication.CreateBuilder(args))
	builder := web.CreateBuilder()

	// Register services
	builder.Services.AddSingleton(NewUserService)

	// Configure Swagger
	swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
		opts.Title = "User API"
		opts.Version = "v1"
		opts.Description = "A simple user management API demonstrating HttpContext and ActionResult"
		opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Enter JWT token",
		})
	})

	// Build application
	app := builder.Build()

	// Configure middleware
	swagger.UseSwagger(app)
	swagger.UseSwaggerUI(app)

	// ==================== Route Examples ====================

	// Style 1: Traditional gin.HandlerFunc (still supported)
	app.MapGet("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to User API",
			"version": "v1",
			"docs":    "/swagger",
		})
	}).WithSummary("Get API information")

	// Style 2: HttpContext with ActionResult (recommended for clean code)
	app.MapGet("/health", func(c *web.HttpContext) web.IActionResult {
		return c.Ok(gin.H{"status": "healthy"})
	}).WithSummary("Health check")

	// User routes with group
	users := app.MapGroup("/api/users")
	users.WithTags("Users")

	// GET /api/users - List all users with ActionResult
	users.MapGet("", func(c *web.HttpContext) web.IActionResult {
		svc := di.GetRequiredService[*UserService](app.Services)
		return c.Ok(svc.ListUsers())
	}).
		WithSummary("List all users").
		WithDescription("Returns a list of all users in the system")

	// GET /api/users/:id - Get user by ID with ActionResult
	users.MapGet("/:id", func(c *web.HttpContext) web.IActionResult {
		svc := di.GetRequiredService[*UserService](app.Services)

		// Use MustPathInt for automatic error handling
		id, errResult := c.MustPathInt("id")
		if errResult != nil {
			return errResult
		}

		user := svc.GetUser(id)
		if user == nil {
			return c.NotFound("User not found")
		}

		return c.Ok(user)
	}).
		WithSummary("Get user by ID").
		WithDescription("Returns a single user by their ID")

	// POST /api/users - Create user with ActionResult
	users.MapPost("", func(c *web.HttpContext) web.IActionResult {
		svc := di.GetRequiredService[*UserService](app.Services)

		var req CreateUserRequest
		if err := c.MustBindJSON(&req); err != nil {
			return err
		}

		user := svc.CreateUser(req.Name, req.Email)
		return c.Created(user)
	}).
		WithSummary("Create a new user").
		WithDescription("Creates a new user with the provided name and email")

	// DELETE /api/users/:id - Delete user with ActionResult
	users.MapDelete("/:id", func(c *web.HttpContext) web.IActionResult {
		svc := di.GetRequiredService[*UserService](app.Services)

		id, errResult := c.MustPathInt("id")
		if errResult != nil {
			return errResult
		}

		if !svc.DeleteUser(id) {
			return c.NotFound("User not found")
		}

		return c.NoContent()
	}).
		WithSummary("Delete a user").
		WithDescription("Deletes a user by their ID")

	// ==================== Additional Examples ====================

	// Example: Using HttpContext without ActionResult
	app.MapGet("/api/time", func(c *web.HttpContext) {
		// Can still use gin.Context methods directly
		c.JSON(200, gin.H{"message": "Use c.JSON directly if preferred"})
	}).WithSummary("Get server time")

	// Example: Using static result helpers
	app.MapGet("/api/redirect-example", func(c *web.HttpContext) web.IActionResult {
		return web.Redirect("/")
	}).WithSummary("Redirect example")

	// Run application
	println("Starting server on http://localhost:8080")
	println("Swagger UI available at http://localhost:8080/swagger")
	println("")
	println("Try these endpoints:")
	println("  GET  /                    - API info (gin.HandlerFunc)")
	println("  GET  /health              - Health check (ActionResult)")
	println("  GET  /api/users           - List users")
	println("  GET  /api/users/1         - Get user by ID")
	println("  POST /api/users           - Create user")
	println("  DELETE /api/users/1       - Delete user")
	app.Run()
}
