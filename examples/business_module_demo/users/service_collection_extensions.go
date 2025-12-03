package users

import (
	"github.com/gocrud/csgo/di"
)

// AddUserServices registers all user-related services.
// This is the extension method pattern from .NET.
//
// Usage:
//
//	builder := web.CreateBuilder()
//	users.AddUserServices(builder.Services)
func AddUserServices(services di.IServiceCollection) {
	// Register repository (Singleton for in-memory store)
	services.AddSingleton(NewUserRepository)

	// Register service (Singleton)
	services.AddSingleton(NewUserService)
}

// AddUserServicesWithOptions registers user services with custom options.
//
// Usage:
//
//	users.AddUserServicesWithOptions(builder.Services, func(opts *UserModuleOptions) {
//	    opts.EnableCache = true
//	    opts.CacheExpiration = 5 * time.Minute
//	})
func AddUserServicesWithOptions(services di.IServiceCollection, configure func(*UserModuleOptions)) {
	opts := &UserModuleOptions{
		EnableCache:     false,
		CacheExpiration: 0,
	}

	if configure != nil {
		configure(opts)
	}

	// Register options
	services.AddSingleton(func() *UserModuleOptions {
		return opts
	})

	// Register services
	AddUserServices(services)
}
