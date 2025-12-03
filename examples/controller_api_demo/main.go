package main

import (
	"controller_api_demo/controllers"
	"controller_api_demo/services"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/openapi"
	"github.com/gocrud/csgo/swagger"
	"github.com/gocrud/csgo/web"
)

func main() {
	// ========================================
	// 1. Create web application builder
	// ========================================
	builder := web.CreateBuilder()

	// ========================================
	// 2. Register services (DI)
	// ========================================
	services.AddServices(builder.Services)
	controllers.AddControllers(builder.Services)

	// ========================================
	// 3. Configure Swagger
	// ========================================
	swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
		opts.Title = "Controller API Demo"
		opts.Version = "v1"
		opts.Description = "Demonstrates ASP.NET Core Controller pattern in Go"
		opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Enter JWT token",
		})
	})

	// ========================================
	// 4. Build application
	// ========================================
	app := builder.Build()

	// ========================================
	// 5. Configure middleware
	// ========================================
	swagger.UseSwagger(app)
	swagger.UseSwaggerUI(app)

	// ========================================
	// 6. Map controllers (similar to app.MapControllers() in .NET)
	// ========================================
	app.MapControllers()

	// ========================================
	// 7. Add root endpoint
	// ========================================
	app.MapGet("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "Controller API Demo",
			"version":     "v1",
			"docs":        "/swagger",
			"controllers": []string{"UserController", "OrderController"},
		})
	}).WithSummary("Get API information")

	// ========================================
	// 8. Run application
	// ========================================
	println("========================================")
	println("Controller API Demo")
	println("========================================")
	println("Server: http://localhost:8080")
	println("Swagger: http://localhost:8080/swagger")
	println("")
	println("Controllers:")
	println("  ✅ UserController  - /api/users")
	println("  ✅ OrderController - /api/orders")
	println("========================================")

	app.Run()
}
