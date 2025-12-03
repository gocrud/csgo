package orders

import (
	"github.com/gocrud/csgo/di"
)

// AddOrderServices registers all order-related services.
// This follows the .NET extension method pattern.
//
// Usage:
//
//	builder := web.CreateBuilder()
//	orders.AddOrderServices(builder.Services)
func AddOrderServices(services di.IServiceCollection) {
	// Register order service as Singleton
	services.AddSingleton(NewOrderService)
}

// AddOrderServicesScoped registers order services with Transient lifetime.
// Use this when you need a new instance per request.
//
// Usage:
//
//	orders.AddOrderServicesScoped(builder.Services)
func AddOrderServicesScoped(services di.IServiceCollection) {
	// Register order service as Transient (per-request in web context)
	services.AddTransient(NewOrderService)
}

