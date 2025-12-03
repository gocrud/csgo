package services

import "github.com/gocrud/csgo/di"

// AddServices registers all application services.
func AddServices(services di.IServiceCollection) {
	// Register services as Singleton
	services.AddSingleton(NewUserService)
	services.AddSingleton(NewOrderService)
}

