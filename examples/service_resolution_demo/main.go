package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// UserService is a simple service for demonstration.
type UserService struct {
	name string
}

func NewUserService() *UserService {
	return &UserService{name: "UserService"}
}

func (s *UserService) GetMessage() string {
	return fmt.Sprintf("Hello from %s", s.name)
}

func main() {
	builder := web.CreateBuilder()

	// Register service
	builder.Services.AddSingleton(NewUserService)

	app := builder.Build()

	// Demo: Style 1 - Traditional way (most explicit)
	app.MapGet("/style1", func(c *gin.Context) {
		var svc *UserService
		app.Services.GetRequiredService(&svc) // ✅ Direct access, with type hints

		c.JSON(200, gin.H{
			"style":   "Traditional",
			"message": svc.GetMessage(),
		})
	})

	// Demo: Style 2 - Generic helper (most concise) ✅ Recommended
	app.MapGet("/style2", func(c *gin.Context) {
		svc := di.GetRequiredService[*UserService](app.Services) // ✅ One line!

		c.JSON(200, gin.H{
			"style":   "Generic Helper",
			"message": svc.GetMessage(),
		})
	})

	// Demo: Style 3 - With error handling
	app.MapGet("/style3", func(c *gin.Context) {
		svc, err := di.GetService[*UserService](app.Services)
		if err != nil {
			c.JSON(500, gin.H{"error": "Service not available"})
			return
		}

		c.JSON(200, gin.H{
			"style":   "With Error Handling",
			"message": svc.GetMessage(),
		})
	})

	fmt.Println("=== Service Resolution Demo ===")
	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("")
	fmt.Println("Try these endpoints:")
	fmt.Println("  - http://localhost:8080/style1  (Traditional)")
	fmt.Println("  - http://localhost:8080/style2  (Generic Helper - Recommended)")
	fmt.Println("  - http://localhost:8080/style3  (With Error Handling)")
	fmt.Println("")

	app.Run()
}
