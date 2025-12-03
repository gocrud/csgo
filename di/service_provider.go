package di

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/gocrud/csgo/di/internal"
)

// IServiceProvider defines the mechanism for retrieving service objects.
type IServiceProvider interface {
	// GetService attempts to retrieve a service and populate it into the target pointer.
	// Example: var svc IUserService; err := provider.GetService(&svc)
	GetService(target interface{}) error

	// GetRequiredService retrieves a service and populates it into the target pointer.
	// Panics if the service cannot be resolved.
	// Example: var svc IUserService; provider.GetRequiredService(&svc)
	GetRequiredService(target interface{})

	// TryGetService attempts to retrieve a service, returns false if not found.
	// Example: var svc IUserService; if provider.TryGetService(&svc) { ... }
	TryGetService(target interface{}) bool

	// GetServices retrieves all services of the specified type.
	// Example: var services []IHostedService; provider.GetServices(&services)
	GetServices(target interface{}) error

	// GetKeyedService retrieves a named service.
	// Example: var svc IUserService; provider.GetKeyedService(&svc, "primary")
	GetKeyedService(target interface{}, serviceKey string) error

	// GetRequiredKeyedService retrieves a required named service.
	// Example: var svc IUserService; provider.GetRequiredKeyedService(&svc, "primary")
	GetRequiredKeyedService(target interface{}, serviceKey string)

	// IsService checks if a service is registered.
	IsService(serviceType reflect.Type) bool

	// Dispose releases all resources.
	Dispose() error
}

// serviceProvider is the concrete implementation of IServiceProvider.
type serviceProvider struct {
	engine   *internal.Engine
	disposed atomic.Bool
}

// GetService retrieves a service and populates it into the target pointer.
func (p *serviceProvider) GetService(target interface{}) error {
	if p.disposed.Load() {
		return errors.New("service provider is disposed")
	}

	// Validate target is a pointer
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	if val.IsNil() {
		return errors.New("target pointer cannot be nil")
	}

	// Get target element type
	elem := val.Elem()
	elemType := elem.Type()

	// Validate type can be set
	if !elem.CanSet() {
		return fmt.Errorf("target of type %v cannot be set", elemType)
	}

	// Resolve service from engine
	service, err := p.engine.Resolve(elemType, "")
	if err != nil {
		return fmt.Errorf("failed to resolve service %v: %w", elemType, err)
	}

	if service == nil {
		return fmt.Errorf("service of type %v is not registered", elemType)
	}

	// Validate type compatibility
	serviceVal := reflect.ValueOf(service)
	if !serviceVal.Type().AssignableTo(elemType) {
		return fmt.Errorf("service type %v is not assignable to target type %v",
			serviceVal.Type(), elemType)
	}

	// Set value
	elem.Set(serviceVal)
	return nil
}

// GetRequiredService retrieves a required service.
func (p *serviceProvider) GetRequiredService(target interface{}) {
	if err := p.GetService(target); err != nil {
		panic(fmt.Sprintf("Failed to resolve required service: %v", err))
	}
}

// TryGetService attempts to retrieve a service.
func (p *serviceProvider) TryGetService(target interface{}) bool {
	err := p.GetService(target)
	return err == nil
}

// GetServices retrieves all services of the specified type.
func (p *serviceProvider) GetServices(target interface{}) error {
	if p.disposed.Load() {
		return errors.New("service provider is disposed")
	}

	// Validate target is a pointer to slice
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer to slice")
	}

	if val.IsNil() {
		return errors.New("target pointer cannot be nil")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Slice {
		return errors.New("target must be a pointer to slice")
	}

	// Get slice element type
	elemType := elem.Type().Elem()

	// Resolve all services from engine
	services, err := p.engine.ResolveAll(elemType)
	if err != nil {
		return fmt.Errorf("failed to resolve services %v: %w", elemType, err)
	}

	// Create slice
	slice := reflect.MakeSlice(elem.Type(), len(services), len(services))
	for i, service := range services {
		slice.Index(i).Set(reflect.ValueOf(service))
	}

	// Set value
	elem.Set(slice)
	return nil
}

// GetKeyedService retrieves a named service.
func (p *serviceProvider) GetKeyedService(target interface{}, serviceKey string) error {
	if p.disposed.Load() {
		return errors.New("service provider is disposed")
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("target must be a non-nil pointer")
	}

	elem := val.Elem()
	elemType := elem.Type()

	if !elem.CanSet() {
		return fmt.Errorf("target of type %v cannot be set", elemType)
	}

	// Resolve named service from engine
	service, err := p.engine.Resolve(elemType, serviceKey)
	if err != nil {
		return fmt.Errorf("failed to resolve keyed service %v[%s]: %w",
			elemType, serviceKey, err)
	}

	if service == nil {
		return fmt.Errorf("keyed service %v[%s] is not registered",
			elemType, serviceKey)
	}

	elem.Set(reflect.ValueOf(service))
	return nil
}

// GetRequiredKeyedService retrieves a required named service.
func (p *serviceProvider) GetRequiredKeyedService(target interface{}, serviceKey string) {
	if err := p.GetKeyedService(target, serviceKey); err != nil {
		panic(fmt.Sprintf("Failed to resolve required keyed service: %v", err))
	}
}

// IsService checks if a service is registered.
func (p *serviceProvider) IsService(serviceType reflect.Type) bool {
	return p.engine.Contains(serviceType, "")
}

// Dispose releases all resources.
func (p *serviceProvider) Dispose() error {
	if !p.disposed.CompareAndSwap(false, true) {
		return nil // Already disposed
	}

	// TODO: Dispose all disposable singletons
	return nil
}
