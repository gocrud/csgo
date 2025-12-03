package web

import (
	"github.com/gocrud/csgo/di"
)

// ⚠️ IMPORTANT: Controller Lifecycle
//
// Controllers are SINGLETONS and created once at application startup.
// They are shared across all HTTP requests for the lifetime of the application.
//
// Best Practices:
// 1. Controllers MUST be stateless - do NOT store request-specific data in controller fields
// 2. Access request data through HttpContext parameter in handlers
// 3. Inject services (IServiceProvider or specific services) via constructor
// 4. For request-scoped services, resolve them using di.GetRequiredService() in handlers
//
// Example:
//
//	type UserController struct {
//	    web.ControllerBase
//	    userService IUserService  // ✅ Service dependency (safe)
//	}
//
//	func (c *UserController) MapRoutes(app *web.WebApplication) {
//	    app.MapGet("/users/:id", func(ctx *web.HttpContext) web.IActionResult {
//	        id, _ := ctx.PathInt("id")
//	        user := c.userService.GetUser(id)  // ✅ Safe: service handles business logic
//	        return ctx.Ok(user)
//	    })
//	}

// IController defines the interface for controllers.
// Controllers that implement this interface can be automatically discovered
// and registered by MapControllers().
type IController interface {
	// MapRoutes registers the controller's routes with the application.
	MapRoutes(app *WebApplication)
}

// ControllerBase provides common functionality for controllers.
// Embed this in your controllers to get access to common services.
type ControllerBase struct {
	Services di.IServiceProvider
}

// NewControllerBase creates a new ControllerBase with the given service provider.
func NewControllerBase(services di.IServiceProvider) ControllerBase {
	return ControllerBase{Services: services}
}

// ControllerOptions represents controller configuration options.
type ControllerOptions struct {
	// EnableEndpointMetadata enables endpoint metadata for OpenAPI generation
	EnableEndpointMetadata bool
}

// controllerRegistry stores registered controller factories
var controllerFactories []func(di.IServiceProvider) IController

// AddControllers adds MVC controller services and enables controller discovery.
// TODO: This method is not used yet.
// Corresponds to .NET services.AddControllers().
func (b *WebApplicationBuilder) AddControllers(configure ...func(*ControllerOptions)) *WebApplicationBuilder {
	opts := &ControllerOptions{
		EnableEndpointMetadata: true,
	}
	if len(configure) > 0 && configure[0] != nil {
		configure[0](opts)
	}

	// Store options for later use
	b.Services.AddSingleton(func() *ControllerOptions { return opts })

	return b
}

// AddController registers a controller factory for automatic discovery.
// The controller will be created as a SINGLETON when MapControllers() is called.
//
// Important: Controllers are singletons and must be stateless. Do not store
// request-specific data in controller fields.
//
// Usage:
//
//	// With constructor function
//	web.AddController(builder.Services, NewUserController)
//
//	// With inline factory
//	web.AddController(builder.Services, func(sp di.IServiceProvider) *UserController {
//	    userService := di.GetRequiredService[IUserService](sp)
//	    return NewUserController(userService)
//	})
func AddController[T IController](services di.IServiceCollection, factory func(di.IServiceProvider) T) {
	// Store factory for later instantiation during MapControllers()
	// Controllers are created once at startup and reused for all requests
	controllerFactories = append(controllerFactories, func(sp di.IServiceProvider) IController {
		return factory(sp)
	})
}

// AddControllerInstance registers an existing controller instance.
// Use this when you need more control over controller creation.
//
// Usage:
//
//	web.AddControllerInstance(builder.Services, func(sp di.IServiceProvider) web.IController {
//	    return NewUserController(sp)
//	})
func AddControllerInstance(services di.IServiceCollection, factory func(di.IServiceProvider) IController) {
	controllerFactories = append(controllerFactories, factory)
}

// MapControllers discovers and registers all controllers as singletons.
// Each controller is created once at startup and used for the lifetime of the application.
// Controllers must be registered using AddController() before calling this method.
//
// This method should be called after Build() and before Run().
// Corresponds to .NET app.MapControllers().
//
// Usage:
//
//	app := builder.Build()
//	app.MapControllers()  // Controllers created here as singletons
//	app.Run()
func (app *WebApplication) MapControllers() *WebApplication {
	// Create each controller once and register its routes
	// These controller instances will be reused for all requests
	for _, factory := range controllerFactories {
		controller := factory(app.Services)
		controller.MapRoutes(app)
	}

	return app
}

// ResetControllers clears all registered controller factories.
// This is mainly useful for testing.
func ResetControllers() {
	controllerFactories = nil
}

// GetRegisteredControllerCount returns the number of registered controllers.
// This is mainly useful for testing and debugging.
func GetRegisteredControllerCount() int {
	return len(controllerFactories)
}
