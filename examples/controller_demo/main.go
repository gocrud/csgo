package main

import (
	"controller_demo/controllers"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/openapi"
	"github.com/gocrud/csgo/swagger"
	"github.com/gocrud/csgo/web"
)

func main() {
	// Create web application builder
	builder := web.CreateBuilder()

	// ========================================
	// Register Services
	// ========================================
	builder.Services.AddSingleton(controllers.NewUserService)
	builder.Services.AddSingleton(controllers.NewProductService)

	// Configure Swagger
	swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
		opts.Title = "Controller Demo API"
		opts.Version = "v1"
		opts.Description = "Demonstrates the Controller pattern in Ego framework"
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
	// Register Controllers (similar to .NET)
	// ========================================

	// Create and register controllers manually
	// (This demo shows the old style - see controller_api_demo for the new style)
	userController := controllers.NewUserController(app)
	userController.MapRoutes(app)

	productController := controllers.NewProductController(app)
	productController.MapRoutes(app)

	// ========================================
	// Root endpoint
	// ========================================
	app.MapGet("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "Controller Demo API",
			"version":     "v1",
			"docs":        "/swagger",
			"controllers": []string{"UserController", "ProductController"},
		})
	}).WithSummary("Get API information")

	// Run application
	println("========================================")
	println("Controller Demo API")
	println("========================================")
	println("Server: http://localhost:8080")
	println("Swagger: http://localhost:8080/swagger")
	println("")
	println("Registered Controllers:")
	println("  ✅ UserController    - /api/users")
	println("  ✅ ProductController - /api/products")
	println("========================================")

	app.Run()
}
