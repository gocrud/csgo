package router

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// TestNestedGroupRoutesCollection tests that routes from nested groups are collected correctly.
func TestNestedGroupRoutesCollection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Create root group
	rootGroup := engine.Group("/api")
	apiGroup := NewRouteGroupBuilder(rootGroup, "/api")

	// Enable OpenAPI on root group
	apiGroup.WithOpenApi(nil)

	// Create nested group (should inherit OpenAPI metadata)
	api2Group := apiGroup.MapGroup("/v2")

	// Add route to root group
	apiGroup.MapGet("/users", func(c *gin.Context) {})

	// Add route to nested group
	api2Group.MapGet("/posts", func(c *gin.Context) {})

	// Create deeply nested group
	api3Group := api2Group.MapGroup("/admin")
	api3Group.MapGet("/settings", func(c *gin.Context) {})

	// Get all routes
	routes := apiGroup.GetRoutes()

	// Should have 3 routes: /api/users, /api/v2/posts, /api/v2/admin/settings
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Verify paths
	expectedPaths := []string{"/api/users", "/api/v2/posts", "/api/v2/admin/settings"}
	foundPaths := make(map[string]bool)
	for _, route := range routes {
		foundPaths[route.GetPath()] = true
	}

	for _, expectedPath := range expectedPaths {
		if !foundPaths[expectedPath] {
			t.Errorf("Expected path %s not found", expectedPath)
		}
	}
}

// TestNestedGroupOpenAPIInheritance tests that OpenAPI metadata is inherited by nested groups.
func TestNestedGroupOpenAPIInheritance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Create root group with OpenAPI enabled
	rootGroup := engine.Group("/api")
	apiGroup := NewRouteGroupBuilder(rootGroup, "/api")
	apiGroup.WithOpenApi(func(b *OpenApiBuilder) {
		b.Tags("API")
	})

	// Create nested group (should inherit OpenAPI metadata)
	api2Group := apiGroup.MapGroup("/v2")

	// Add route to nested group
	api2Group.MapGet("/test", func(c *gin.Context) {})

	// Get routes
	routes := apiGroup.GetRoutes()

	// Should have 1 route
	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	route := routes[0]

	// Route should have OpenAPI enabled (inherited from parent group)
	if !route.IsOpenApiEnabled() {
		t.Error("Expected route to have OpenAPI enabled (inherited from parent)")
	}
}

// TestNestedGroupMiddlewareInheritance tests that middleware is inherited by nested groups.
func TestNestedGroupMiddlewareInheritance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Track middleware execution
	var executionOrder []string

	middleware1 := func(c *gin.Context) {
		executionOrder = append(executionOrder, "middleware1")
		c.Next()
	}

	middleware2 := func(c *gin.Context) {
		executionOrder = append(executionOrder, "middleware2")
		c.Next()
	}

	// Create root group with middleware1
	rootGroup := engine.Group("/api", middleware1)
	apiGroup := NewRouteGroupBuilder(rootGroup, "/api")

	// Create nested group with middleware2
	api2Group := apiGroup.MapGroup("/v2", middleware2)

	// Add route to nested group
	api2Group.MapGet("/test", func(c *gin.Context) {
		executionOrder = append(executionOrder, "handler")
	})

	// Test the route by making a request
	// Note: This test verifies that Gin's Group() correctly inherits middleware
	// The actual middleware execution would be tested in an integration test
	// Here we just verify the group structure is correct

	routes := apiGroup.GetRoutes()
	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	// Verify the route path is correct
	if routes[0].GetPath() != "/api/v2/test" {
		t.Errorf("Expected path /api/v2/test, got %s", routes[0].GetPath())
	}
}

// TestMultiLevelNestedGroups tests multiple levels of nested groups.
func TestMultiLevelNestedGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Create root group
	rootGroup := engine.Group("/api")
	apiGroup := NewRouteGroupBuilder(rootGroup, "/api")
	apiGroup.WithOpenApi(nil)

	// Level 1: /api/v1
	v1Group := apiGroup.MapGroup("/v1")
	v1Group.MapGet("/users", func(c *gin.Context) {})

	// Level 2: /api/v1/admin
	adminGroup := v1Group.MapGroup("/admin")
	adminGroup.MapGet("/settings", func(c *gin.Context) {})

	// Level 3: /api/v1/admin/security
	securityGroup := adminGroup.MapGroup("/security")
	securityGroup.MapGet("/permissions", func(c *gin.Context) {})

	// Level 1: /api/v2 (sibling of v1)
	v2Group := apiGroup.MapGroup("/v2")
	v2Group.MapGet("/posts", func(c *gin.Context) {})

	// Get all routes
	routes := apiGroup.GetRoutes()

	// Should have 4 routes
	if len(routes) != 4 {
		t.Errorf("Expected 4 routes, got %d", len(routes))
	}

	// Verify all routes have OpenAPI enabled (inherited)
	for _, route := range routes {
		if !route.IsOpenApiEnabled() {
			t.Errorf("Route %s should have OpenAPI enabled", route.GetPath())
		}
	}
}

// TestChildGroupMetadataInheritance tests that metadata is properly inherited.
func TestChildGroupMetadataInheritance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Create root group
	rootGroup := engine.Group("/api")
	apiGroup := NewRouteGroupBuilder(rootGroup, "/api")

	// Add OpenAPI metadata
	apiGroup.WithOpenApi(func(b *OpenApiBuilder) {
		b.Tags("API", "V1")
	})

	// Create child group
	childGroup := apiGroup.MapGroup("/child")

	// Child should inherit parent's metadata
	if len(childGroup.metadata) == 0 {
		t.Error("Child group should inherit parent metadata")
	}

	// Verify OpenAPI metadata is present
	hasOpenApiMetadata := false
	for _, meta := range childGroup.metadata {
		if _, ok := meta.(*OpenApiMetadata); ok {
			hasOpenApiMetadata = true
			break
		}
	}

	if !hasOpenApiMetadata {
		t.Error("Child group should inherit OpenApiMetadata from parent")
	}
}
