package main

import (
	"strconv"

	"business_module_demo/orders"
	"business_module_demo/users"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/openapi"
	"github.com/gocrud/csgo/swagger"
	"github.com/gocrud/csgo/web"
)

func main() {
	// Create web application builder
	builder := web.CreateBuilder()

	// ========================================
	// Style 1: Simple registration (most common)
	// ========================================
	users.AddUserServices(builder.Services)
	orders.AddOrderServices(builder.Services)

	// ========================================
	// Style 2: Registration with options
	// ========================================
	// users.AddUserServicesWithOptions(builder.Services, func(opts *users.UserModuleOptions) {
	// 	opts.EnableCache = true
	// 	opts.CacheExpiration = 5 * time.Minute
	// })

	// ========================================
	// Style 3: Scoped registration (Transient)
	// ========================================
	// orders.AddOrderServicesScoped(builder.Services)

	// Configure Swagger
	swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
		opts.Title = "Business Module Demo API"
		opts.Version = "v1"
		opts.Description = "Demonstrates how to create IServiceCollection extensions for business modules"
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

	// ========================================
	// User API Routes
	// ========================================
	userGroup := app.MapGroup("/api/users")
	userGroup.WithTags("Users")

	// GET /api/users
	userGroup.MapGet("", func(c *gin.Context) {
		// ✅ Style 1: Traditional way
		var userSvc users.IUserService
		app.Services.GetRequiredService(&userSvc)

		userList, err := userSvc.ListUsers()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, userList)
	}).
		WithSummary("List all users").
		WithDescription("Returns a list of all users in the system")

	// GET /api/users/{id}
	userGroup.MapGet("/{id}", func(c *gin.Context) {
		// ✅ Style 2: Generic helper (recommended)
		userSvc := di.GetRequiredService[users.IUserService](app.Services)

		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		user, err := userSvc.GetUser(idInt)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, user)
	}).
		WithSummary("Get user by ID").
		WithDescription("Returns a single user by their ID")

	// POST /api/users
	userGroup.MapPost("", func(c *gin.Context) {
		userSvc := di.GetRequiredService[users.IUserService](app.Services)

		var user users.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := userSvc.CreateUser(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, user)
	}).
		WithSummary("Create a new user").
		WithDescription("Creates a new user in the system")

	// ========================================
	// Order API Routes
	// ========================================
	orderGroup := app.MapGroup("/api/orders")
	orderGroup.WithTags("Orders")

	// GET /api/orders/{id}
	orderGroup.MapGet("/{id}", func(c *gin.Context) {
		orderSvc := di.GetRequiredService[orders.OrderService](app.Services)

		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		order, err := orderSvc.GetOrder(idInt)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, order)
	}).
		WithSummary("Get order by ID").
		WithDescription("Returns a single order by its ID")

	// GET /api/orders/user/{userId}
	orderGroup.MapGet("/user/{userId}", func(c *gin.Context) {
		orderSvc := di.GetRequiredService[orders.OrderService](app.Services)

		userID := c.Param("userId")
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}

		orderList, err := orderSvc.ListOrdersByUser(userIDInt)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, orderList)
	}).
		WithSummary("List orders by user").
		WithDescription("Returns all orders for a specific user")

	// POST /api/orders
	orderGroup.MapPost("", func(c *gin.Context) {
		orderSvc := di.GetRequiredService[orders.OrderService](app.Services)

		var order orders.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := orderSvc.CreateOrder(&order); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, order)
	}).
		WithSummary("Create a new order").
		WithDescription("Creates a new order in the system")

	// ========================================
	// Root endpoint
	// ========================================
	app.MapGet("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Business Module Demo API",
			"version": "v1",
			"docs":    "/swagger",
			"modules": []string{"users", "orders"},
		})
	}).WithSummary("Get API information")

	// Run application
	println("========================================")
	println("Business Module Demo API")
	println("========================================")
	println("Server: http://localhost:8080")
	println("Swagger: http://localhost:8080/swagger")
	println("")
	println("Registered Modules:")
	println("  ✅ Users Module   - /api/users")
	println("  ✅ Orders Module  - /api/orders")
	println("========================================")

	app.Run()
}
