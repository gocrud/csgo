package web

import (
	"github.com/gocrud/csgo/di"
)

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

// AddController registers a controller with the DI container.
// The controller will be automatically discovered and registered when MapControllers() is called.
//
// Usage:
//
//	web.AddController[*UserController](builder.Services, NewUserController)
func AddController[T IController](services di.IServiceCollection, factory func(di.IServiceProvider) T) {
	// Register the concrete controller type as scoped
	services.AddScoped(func(sp di.IServiceProvider) T {
		return factory(sp)
	})

	// Register a factory that creates IController from this controller
	// This allows MapControllers to discover all controllers
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

// MapControllers discovers and registers all controllers.
// Controllers must be registered using AddController() before calling this method.
// Corresponds to .NET app.MapControllers().
//
// Usage:
//
//	app := builder.Build()
//	app.MapControllers()
//	app.Run()
func (app *WebApplication) MapControllers() *WebApplication {
	// Create and register routes for each controller
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
