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
	// Get retrieves a service and populates it into the target pointer (panics if not found).
	// Supports both pointer and value types with automatic dereferencing:
	//   - var svc *UserService; provider.Get(&svc)  // pointer type (zero-copy)
	//   - var svc UserService; provider.Get(&svc)   // value type (auto-deref + copy)
	Get(target interface{})

	// GetNamed retrieves a named service and populates it into the target pointer (panics if not found).
	//   - var db *Database; provider.GetNamed(&db, "primary")
	GetNamed(target interface{}, serviceKey string)

	// Internal methods (used by generic API functions)
	resolveType(t reflect.Type) (interface{}, error)
	resolveNamed(t reflect.Type, name string) (interface{}, error)
	resolveAll(t reflect.Type) []interface{}

	// Dispose releases all resources.
	Dispose() error
}

// serviceProvider is the concrete implementation of IServiceProvider.
type serviceProvider struct {
	engine   *internal.Engine
	disposed atomic.Bool
}

// Get retrieves a service and populates it into the target pointer.
// Supports both pointer and value types with intelligent dereferencing.
func (p *serviceProvider) Get(target interface{}) {
	if p.disposed.Load() {
		panic("service provider is disposed")
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		panic("target must be a non-nil pointer")
	}

	elem := val.Elem()
	elemType := elem.Type()

	// Try 1: Direct lookup for the target type
	service, err := p.engine.Resolve(elemType, "")
	if err == nil {
		elem.Set(reflect.ValueOf(service))
		return
	}

	// Try 2: If target is a value type (struct), try to find pointer type and auto-deref
	if elemType.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(elemType)
		ptrService, ptrErr := p.engine.Resolve(ptrType, "")
		if ptrErr == nil {
			// Auto-dereference: assign copy of the value
			elem.Set(reflect.ValueOf(ptrService).Elem())
			return
		}
	}

	panic(fmt.Sprintf("service %v not found", elemType))
}

// GetNamed retrieves a named service and populates it into the target pointer.
func (p *serviceProvider) GetNamed(target interface{}, serviceKey string) {
	if p.disposed.Load() {
		panic("service provider is disposed")
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		panic("target must be a non-nil pointer")
	}

	elem := val.Elem()
	elemType := elem.Type()

	// Try 1: Direct lookup
	service, err := p.engine.Resolve(elemType, serviceKey)
	if err == nil {
		elem.Set(reflect.ValueOf(service))
		return
	}

	// Try 2: Auto-deref for value types
	if elemType.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(elemType)
		ptrService, ptrErr := p.engine.Resolve(ptrType, serviceKey)
		if ptrErr == nil {
			elem.Set(reflect.ValueOf(ptrService).Elem())
			return
		}
	}

	panic(fmt.Sprintf("named service %v[%s] not found", elemType, serviceKey))
}

// resolveType resolves a service by type (internal method for generic API).
func (p *serviceProvider) resolveType(t reflect.Type) (interface{}, error) {
	if p.disposed.Load() {
		return nil, errors.New("provider disposed")
	}
	return p.engine.Resolve(t, "")
}

// resolveNamed resolves a named service by type (internal method for generic API).
func (p *serviceProvider) resolveNamed(t reflect.Type, name string) (interface{}, error) {
	if p.disposed.Load() {
		return nil, errors.New("provider disposed")
	}
	return p.engine.Resolve(t, name)
}

// resolveAll resolves all services of a type (internal method for generic API).
func (p *serviceProvider) resolveAll(t reflect.Type) []interface{} {
	if p.disposed.Load() {
		return nil
	}
	services, _ := p.engine.ResolveAll(t)
	return services
}

// Dispose releases all resources including singleton services.
// All singleton services implementing IDisposable will have their Dispose method called.
func (p *serviceProvider) Dispose() error {
	if !p.disposed.CompareAndSwap(false, true) {
		return nil // Already disposed
	}

	// Dispose all singleton services implementing IDisposable
	singletons := p.engine.GetSingletons()
	var errors []error

	// Dispose in reverse order (LIFO)
	for i := len(singletons) - 1; i >= 0; i-- {
		if disposable, ok := singletons[i].(IDisposable); ok {
			if err := disposable.Dispose(); err != nil {
				errors = append(errors, fmt.Errorf("failed to dispose singleton at index %d: %w", i, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("provider disposal encountered %d error(s): %v", len(errors), errors)
	}

	return nil
}
