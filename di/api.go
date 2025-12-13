package di

import (
	"fmt"
	"reflect"
)

// Get retrieves a service from the provider (panics if not found).
// Supports both pointer and value types with auto-dereferencing:
//   - userSvc := di.Get[*UserService](provider)  // pointer type (zero-copy)
//   - config := di.Get[AppConfig](provider)      // value type (auto-deref + copy)
func Get[T any](provider IServiceProvider) T {
	var zero T
	// Use reflect.TypeOf with pointer trick to handle interface types correctly
	t := reflect.TypeOf((*T)(nil)).Elem()

	// Try 1: Direct lookup
	service, err := provider.resolveType(t)
	if err == nil {
		return service.(T)
	}

	// Try 2: If T is a value type (struct), try to find pointer type and auto-deref
	if t.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(t)
		ptrService, ptrErr := provider.resolveType(ptrType)
		if ptrErr == nil {
			// Auto-dereference: return copy of the value
			return reflect.ValueOf(ptrService).Elem().Interface().(T)
		}
	}

	panic(fmt.Sprintf("service %T not found: %v", zero, err))
}

// GetOr retrieves a service from the provider, returns defaultValue if not found.
// Does not panic, useful for optional services. Supports auto-dereferencing.
//   - config := di.GetOr[*Config](provider, defaultConfig)
//   - val := di.GetOr[AppConfig](provider, defaultVal)
func GetOr[T any](provider IServiceProvider, defaultValue T) T {
	// Use reflect.TypeOf with pointer trick to handle interface types correctly
	t := reflect.TypeOf((*T)(nil)).Elem()

	// Try 1: Direct lookup
	service, err := provider.resolveType(t)
	if err == nil {
		return service.(T)
	}

	// Try 2: If T is a value type (struct), try to find pointer type and auto-deref
	if t.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(t)
		ptrService, ptrErr := provider.resolveType(ptrType)
		if ptrErr == nil {
			// Auto-dereference
			return reflect.ValueOf(ptrService).Elem().Interface().(T)
		}
	}

	return defaultValue
}

// GetAll retrieves all services of the specified type.
// Returns a slice of all registered services matching the type.
//   - plugins := di.GetAll[IPlugin](provider)
func GetAll[T any](provider IServiceProvider) []T {
	// Use reflect.TypeOf with pointer trick to handle interface types correctly
	t := reflect.TypeOf((*T)(nil)).Elem()

	services := provider.resolveAll(t)
	result := make([]T, len(services))
	for i, svc := range services {
		result[i] = svc.(T)
	}
	return result
}

// GetNamed retrieves a named service from the provider (panics if not found).
// Useful for services registered with a specific key.
// Supports auto-dereferencing like Get:
//   - primaryDB := di.GetNamed[*Database](provider, "primary")  // pointer
//   - config := di.GetNamed[AppConfig](provider, "app")         // value (auto-deref)
func GetNamed[T any](provider IServiceProvider, name string) T {
	var zero T
	// Use reflect.TypeOf with pointer trick to handle interface types correctly
	t := reflect.TypeOf((*T)(nil)).Elem()

	// Try 1: Direct lookup
	service, err := provider.resolveNamed(t, name)
	if err == nil {
		return service.(T)
	}

	// Try 2: If T is a value type (struct), try to find pointer type and auto-deref
	if t.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(t)
		ptrService, ptrErr := provider.resolveNamed(ptrType, name)
		if ptrErr == nil {
			// Auto-dereference
			return reflect.ValueOf(ptrService).Elem().Interface().(T)
		}
	}

	panic(fmt.Sprintf("named service %T[%s] not found: %v", zero, name, err))
}

// TryGet attempts to retrieve a service from the provider.
// Returns the service and true if found, zero value and false otherwise.
// Supports auto-dereferencing like Get.
//   - if svc, ok := di.TryGet[*Service](provider); ok { ... }
//   - if val, ok := di.TryGet[AppConfig](provider); ok { ... }
func TryGet[T any](provider IServiceProvider) (T, bool) {
	var zero T
	// Use reflect.TypeOf with pointer trick to handle interface types correctly
	t := reflect.TypeOf((*T)(nil)).Elem()

	// Try 1: Direct lookup
	service, err := provider.resolveType(t)
	if err == nil {
		return service.(T), true
	}

	// Try 2: If T is a value type (struct), try to find pointer type and auto-deref
	if t.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(t)
		ptrService, ptrErr := provider.resolveType(ptrType)
		if ptrErr == nil {
			// Auto-dereference
			return reflect.ValueOf(ptrService).Elem().Interface().(T), true
		}
	}

	return zero, false
}
