package controllers

import (
	svc "controller_api_demo/services"

	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// AddControllers registers all controllers with the DI container.
// Controllers will be automatically discovered when app.MapControllers() is called.
func AddControllers(services di.IServiceCollection) {
	// Register controllers using the new web.AddController API
	// This allows automatic discovery by app.MapControllers()
	web.AddController(services, func(sp di.IServiceProvider) *UserController {
		userService := di.GetRequiredService[svc.UserService](sp)
		return NewUserController(userService)
	})

	web.AddController(services, func(sp di.IServiceProvider) *OrderController {
		orderService := di.GetRequiredService[svc.OrderService](sp)
		return NewOrderController(orderService)
	})
}
