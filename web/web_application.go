package web

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/hosting"
	"github.com/gocrud/csgo/web/routing"
)

// WebApplication represents the configured web application.
type WebApplication struct {
	host        hosting.IHost
	engine      *gin.Engine
	Services    di.IServiceProvider // ✅ 直接暴露，强类型
	Environment hosting.IHostEnvironment
	routes      []*routing.RouteBuilder
	groups      []*routing.RouteGroupBuilder
	runtimeUrls *[]string // Pointer to runtime URLs (shared with HttpServer)

	// Handler converter with services injection
	toHandler  func(Handler) gin.HandlerFunc
	toHandlers func(...Handler) []gin.HandlerFunc
}

// Run runs the web application and blocks until shutdown.
// If urls are provided, they override the configured listen addresses.
// Corresponds to .NET app.Run(url).
func (app *WebApplication) Run(urls ...string) error {
	if len(urls) > 0 && app.runtimeUrls != nil {
		*app.runtimeUrls = urls
	}
	return app.host.Run()
}

// RunWithContext runs the web application with a custom context and blocks until shutdown.
func (app *WebApplication) RunWithContext(ctx context.Context) error {
	return app.host.RunWithContext(ctx)
}

// Start starts the web application.
func (app *WebApplication) Start(ctx context.Context) error {
	return app.host.Start(ctx)
}

// Stop stops the web application.
func (app *WebApplication) Stop(ctx context.Context) error {
	return app.host.Stop(ctx)
}

// Use adds middleware to the pipeline.
func (app *WebApplication) Use(middleware ...gin.HandlerFunc) {
	app.engine.Use(middleware...)
}

// MapGet registers a GET endpoint.
// Supports multiple handler types:
//   - gin.HandlerFunc
//   - func(*HttpContext)
//   - func(*HttpContext) IActionResult
func (app *WebApplication) MapGet(pattern string, handlers ...Handler) routing.IEndpointConventionBuilder {
	return app.mapRoute("GET", pattern, handlers...)
}

// MapPost registers a POST endpoint.
// Supports multiple handler types:
//   - gin.HandlerFunc
//   - func(*HttpContext)
//   - func(*HttpContext) IActionResult
func (app *WebApplication) MapPost(pattern string, handlers ...Handler) routing.IEndpointConventionBuilder {
	return app.mapRoute("POST", pattern, handlers...)
}

// MapPut registers a PUT endpoint.
// Supports multiple handler types:
//   - gin.HandlerFunc
//   - func(*HttpContext)
//   - func(*HttpContext) IActionResult
func (app *WebApplication) MapPut(pattern string, handlers ...Handler) routing.IEndpointConventionBuilder {
	return app.mapRoute("PUT", pattern, handlers...)
}

// MapDelete registers a DELETE endpoint.
// Supports multiple handler types:
//   - gin.HandlerFunc
//   - func(*HttpContext)
//   - func(*HttpContext) IActionResult
func (app *WebApplication) MapDelete(pattern string, handlers ...Handler) routing.IEndpointConventionBuilder {
	return app.mapRoute("DELETE", pattern, handlers...)
}

// MapPatch registers a PATCH endpoint.
// Supports multiple handler types:
//   - gin.HandlerFunc
//   - func(*HttpContext)
//   - func(*HttpContext) IActionResult
func (app *WebApplication) MapPatch(pattern string, handlers ...Handler) routing.IEndpointConventionBuilder {
	return app.mapRoute("PATCH", pattern, handlers...)
}

// MapGroup creates a route group.
// Supports multiple handler types for middleware:
//   - gin.HandlerFunc
//   - func(*HttpContext)
//   - func(*HttpContext) IActionResult
func (app *WebApplication) MapGroup(prefix string, handlers ...Handler) *routing.RouteGroupBuilder {
	// Convert handlers for group middleware using services-aware converter
	ginHandlers := app.toHandlers(handlers...)

	ginGroup := app.engine.Group(prefix, ginHandlers...)
	group := routing.NewRouteGroupBuilder(ginGroup, prefix)

	// Set handler converter with services injection
	group.SetHandlerConverter(app.toHandler)

	app.groups = append(app.groups, group)
	return group
}

// mapRoute is the internal method to register a route.
func (app *WebApplication) mapRoute(method, pattern string, handlers ...Handler) routing.IEndpointConventionBuilder {
	// Convert handlers to gin.HandlerFunc using services-aware converter
	ginHandlers := app.toHandlers(handlers...)

	// Register with Gin
	app.engine.Handle(method, pattern, ginHandlers...)

	// Create route builder
	rb := routing.NewRouteBuilder(method, pattern)
	app.routes = append(app.routes, rb)

	return rb
}

// GetRoutes returns all registered routes.
func (app *WebApplication) GetRoutes() []*routing.RouteBuilder {
	allRoutes := make([]*routing.RouteBuilder, 0)

	// Add top-level routes
	allRoutes = append(allRoutes, app.routes...)

	// Add routes from groups
	for _, group := range app.groups {
		allRoutes = append(allRoutes, group.GetRoutes()...)
	}

	return allRoutes
}
