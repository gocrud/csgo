package di

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/gocrud/csgo/di/internal"
)

// IServiceScope provides a scope for service resolution.
type IServiceScope interface {
	// ServiceProvider returns the service provider for this scope.
	ServiceProvider() IServiceProvider

	// Dispose releases all scoped resources.
	Dispose() error
}

// serviceScope is the concrete implementation of IServiceScope.
type serviceScope struct {
	parent    *serviceProvider
	instances map[internal.RegistrationKey]interface{}
	mu        sync.RWMutex
	disposed  atomic.Bool
}

// ServiceProvider returns the scoped service provider.
func (s *serviceScope) ServiceProvider() IServiceProvider {
	return &scopedServiceProvider{
		scope:  s,
		parent: s.parent,
	}
}

// Dispose releases all scoped resources.
func (s *serviceScope) Dispose() error {
	if !s.disposed.CompareAndSwap(false, true) {
		return nil // Already disposed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Dispose all scoped instances that implement disposable interface
	for _, instance := range s.instances {
		if disposable, ok := instance.(interface{ Dispose() error }); ok {
			disposable.Dispose()
		}
	}

	s.instances = nil
	return nil
}

// scopedServiceProvider is a service provider that resolves services within a scope.
type scopedServiceProvider struct {
	scope  *serviceScope
	parent *serviceProvider
}

// GetService retrieves a service within the scope.
func (sp *scopedServiceProvider) GetService(target interface{}) error {
	if sp.scope.disposed.Load() {
		return errors.New("scope is disposed")
	}

	// Validate target is a pointer
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("target must be a non-nil pointer")
	}

	elem := val.Elem()
	elemType := elem.Type()

	if !elem.CanSet() {
		return fmt.Errorf("target of type %v cannot be set", elemType)
	}

	key := internal.RegistrationKey{Type: elemType, Name: ""}

	// Get registration info
	reg, exists := sp.parent.engine.GetRegistration(key)
	if !exists {
		return fmt.Errorf("service %v not found", elemType)
	}

	var instance interface{}
	var err error

	switch reg.Lifetime {
	case internal.Singleton:
		// Singleton from root container
		return sp.parent.GetService(target)

	case internal.Scoped:
		// Scoped from current scope - get or create
		sp.scope.mu.RLock()
		instance, exists = sp.scope.instances[key]
		sp.scope.mu.RUnlock()

		if exists {
			elem.Set(reflect.ValueOf(instance))
			return nil
		}

		sp.scope.mu.Lock()
		defer sp.scope.mu.Unlock()

		// Double-check
		if instance, exists = sp.scope.instances[key]; exists {
			elem.Set(reflect.ValueOf(instance))
			return nil
		}

		// Create scoped instance
		instance, err = sp.parent.engine.Resolve(reg.ServiceType, "")
		if err != nil {
			return err
		}

		sp.scope.instances[key] = instance
		elem.Set(reflect.ValueOf(instance))
		return nil

	case internal.Transient:
		// Transient - create new instance each time
		instance, err = sp.parent.engine.Resolve(reg.ServiceType, "")
		if err != nil {
			return err
		}
		elem.Set(reflect.ValueOf(instance))
		return nil
	}

	return fmt.Errorf("unknown lifetime for service %v", elemType)
}

// GetRequiredService retrieves a required service within the scope.
func (sp *scopedServiceProvider) GetRequiredService(target interface{}) {
	if err := sp.GetService(target); err != nil {
		panic(fmt.Sprintf("Failed to resolve required service: %v", err))
	}
}

// TryGetService attempts to retrieve a service within the scope.
func (sp *scopedServiceProvider) TryGetService(target interface{}) bool {
	return sp.GetService(target) == nil
}

// GetServices retrieves all services of the specified type within the scope.
func (sp *scopedServiceProvider) GetServices(target interface{}) error {
	// For now, delegate to parent
	return sp.parent.GetServices(target)
}

// GetKeyedService retrieves a named service within the scope.
func (sp *scopedServiceProvider) GetKeyedService(target interface{}, serviceKey string) error {
	// Similar logic to GetService but with serviceKey
	return sp.parent.GetKeyedService(target, serviceKey)
}

// GetRequiredKeyedService retrieves a required named service within the scope.
func (sp *scopedServiceProvider) GetRequiredKeyedService(target interface{}, serviceKey string) {
	if err := sp.GetKeyedService(target, serviceKey); err != nil {
		panic(fmt.Sprintf("Failed to resolve required keyed service: %v", err))
	}
}

// IsService checks if a service is registered.
func (sp *scopedServiceProvider) IsService(serviceType reflect.Type) bool {
	return sp.parent.IsService(serviceType)
}

// CreateScope creates a nested scope.
func (sp *scopedServiceProvider) CreateScope() IServiceScope {
	// Create a new scope from the parent provider
	return sp.parent.CreateScope()
}

// Dispose disposes the scoped provider (delegates to scope).
func (sp *scopedServiceProvider) Dispose() error {
	return sp.scope.Dispose()
}

