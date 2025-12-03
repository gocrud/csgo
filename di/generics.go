package di

import "reflect"

// GetService is a generic helper to retrieve a service.
// Example: svc, err := GetService[IUserService](provider)
func GetService[T any](provider IServiceProvider) (T, error) {
	var result T
	err := provider.GetService(&result)
	return result, err
}

// GetRequiredService is a generic helper to retrieve a required service.
// Example: svc := GetRequiredService[IUserService](provider)
func GetRequiredService[T any](provider IServiceProvider) T {
	var result T
	provider.GetRequiredService(&result)
	return result
}

// GetServices is a generic helper to retrieve all services of a type.
// Example: services, err := GetServices[IHostedService](provider)
func GetServices[T any](provider IServiceProvider) ([]T, error) {
	var results []T
	err := provider.GetServices(&results)
	return results, err
}

// GetKeyedService is a generic helper to retrieve a named service.
// Example: svc, err := GetKeyedService[IUserService](provider, "primary")
func GetKeyedService[T any](provider IServiceProvider, serviceKey string) (T, error) {
	var result T
	err := provider.GetKeyedService(&result, serviceKey)
	return result, err
}

// GetRequiredKeyedService is a generic helper to retrieve a required named service.
// Example: svc := GetRequiredKeyedService[IUserService](provider, "primary")
func GetRequiredKeyedService[T any](provider IServiceProvider, serviceKey string) T {
	var result T
	provider.GetRequiredKeyedService(&result, serviceKey)
	return result
}

// AddSingleton is a generic helper to add a singleton service.
func AddSingleton[TService any](services IServiceCollection, factory func() TService) IServiceCollection {
	return services.AddSingleton(factory)
}

// AddScoped is a generic helper to add a scoped service.
func AddScoped[TService any](services IServiceCollection, factory func() TService) IServiceCollection {
	return services.AddScoped(factory)
}

// AddTransient is a generic helper to add a transient service.
func AddTransient[TService any](services IServiceCollection, factory func() TService) IServiceCollection {
	return services.AddTransient(factory)
}

// AddKeyedSingleton is a generic helper to add a keyed singleton service.
func AddKeyedSingleton[TService any](services IServiceCollection, serviceKey string, factory func() TService) IServiceCollection {
	return services.AddKeyedSingleton(serviceKey, factory)
}

// AddKeyedScoped is a generic helper to add a keyed scoped service.
func AddKeyedScoped[TService any](services IServiceCollection, serviceKey string, factory func() TService) IServiceCollection {
	return services.AddKeyedScoped(serviceKey, factory)
}

// AddKeyedTransient is a generic helper to add a keyed transient service.
func AddKeyedTransient[TService any](services IServiceCollection, serviceKey string, factory func() TService) IServiceCollection {
	return services.AddKeyedTransient(serviceKey, factory)
}

// TypeOf returns the reflect.Type of the generic type T.
// Example: typ := TypeOf[IUserService]()
func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// GetServiceScopeFactory retrieves the IServiceScopeFactory from the provider.
func GetServiceScopeFactory(provider IServiceProvider) IServiceScopeFactory {
	var factory IServiceScopeFactory
	provider.GetRequiredService(&factory)
	return factory
}

