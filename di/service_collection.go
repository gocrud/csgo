package di

import (
	"fmt"
	"reflect"

	"github.com/gocrud/csgo/di/internal"
)

// IServiceCollection is a contract for a collection of service descriptors.
// This is a pure registration interface following the Interface Segregation Principle.
// Build, Count, and GetDescriptors methods are available on the concrete type but not in the interface.
type IServiceCollection interface {
	// Add registers a singleton service using a constructor function.
	Add(constructor interface{}) IServiceCollection

	// AddInstance registers a singleton instance (pre-created object).
	AddInstance(instance interface{}) IServiceCollection

	// AddNamed registers a named singleton service.
	AddNamed(name string, constructor interface{}) IServiceCollection

	// TryAdd attempts to add a singleton service if it doesn't exist.
	TryAdd(constructor interface{}) IServiceCollection

	// AddHostedService registers a hosted service (background service).
	AddHostedService(constructor interface{}) IServiceCollection
}

// serviceCollection is the concrete implementation of IServiceCollection.
type serviceCollection struct {
	engine *internal.Engine
}

// NewServiceCollection creates a new service collection.
func NewServiceCollection() IServiceCollection {
	return &serviceCollection{
		engine: internal.NewEngine(),
	}
}

// BuildServiceProvider builds the service provider from a service collection.
// This function allows building from the interface type.
// Usage: provider := di.BuildServiceProvider(services)
func BuildServiceProvider(services IServiceCollection) IServiceProvider {
	// Type assert to concrete implementation
	sc, ok := services.(*serviceCollection)
	if !ok {
		panic("services must be created by NewServiceCollection")
	}

	return sc.Build()
}

// Add registers a singleton service using a constructor function.
func (s *serviceCollection) Add(constructor interface{}) IServiceCollection {
	if err := s.register(constructor, Singleton); err != nil {
		panic(fmt.Sprintf("failed to register service: %v", err))
	}
	return s
}

// TryAdd attempts to add a singleton if it doesn't exist.
func (s *serviceCollection) TryAdd(constructor interface{}) IServiceCollection {
	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return s
	}
	if ctorType.NumOut() == 0 {
		return s
	}

	returnType := ctorType.Out(0)
	if !s.engine.Contains(returnType, "") {
		s.Add(constructor)
	}
	return s
}

// AddInstance registers a singleton instance (pre-created object).
func (s *serviceCollection) AddInstance(instance interface{}) IServiceCollection {
	if instance == nil {
		panic("instance cannot be nil")
	}

	instanceType := reflect.TypeOf(instance)

	// Create a properly typed constructor using reflection
	// This ensures Register() correctly extracts the type
	factoryType := reflect.FuncOf([]reflect.Type{}, []reflect.Type{instanceType}, false)
	factoryValue := reflect.MakeFunc(factoryType, func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(instance)}
	})

	reg := &internal.Registration{
		ServiceType:        instanceType,
		ImplementationType: instanceType,
		Lifetime:           internal.Singleton,
		Factory:            factoryValue.Interface(),
	}

	if err := s.engine.Register(reg); err != nil {
		panic(fmt.Sprintf("failed to register instance: %v", err))
	}

	return s
}

// AddHostedService registers a hosted service.
// The service will be started when the host starts and stopped when the host stops.
func (s *serviceCollection) AddHostedService(constructor interface{}) IServiceCollection {
	// Register as Singleton (hosted services should be singletons)
	return s.Add(constructor)
}

// AddNamed registers a named singleton service.
func (s *serviceCollection) AddNamed(name string, constructor interface{}) IServiceCollection {
	if err := s.registerKeyed(constructor, Singleton, name); err != nil {
		panic(fmt.Sprintf("failed to register named service: %v", err))
	}
	return s
}

// Build builds the service provider.
// This is a convenience method on the concrete type (not in the interface).
// Usage: provider := services.Build()
func (s *serviceCollection) Build() IServiceProvider {
	provider := &serviceProvider{
		engine: s.engine,
	}

	if err := s.engine.Compile(); err != nil {
		panic(fmt.Sprintf("failed to build service provider: %v", err))
	}

	return provider
}

// Count returns the number of registered services.
// This is a diagnostic method on the concrete type (not in the interface).
func (s *serviceCollection) Count() int {
	registrations := s.engine.GetAllRegistrations()
	return len(registrations)
}

// GetDescriptors returns all service descriptors.
// This is a diagnostic method on the concrete type (not in the interface).
func (s *serviceCollection) GetDescriptors() []ServiceDescriptor {
	registrations := s.engine.GetAllRegistrations()
	descriptors := make([]ServiceDescriptor, 0, len(registrations))

	for key, reg := range registrations {
		descriptor := ServiceDescriptor{
			ServiceType:        reg.ServiceType,
			ImplementationType: reg.ImplementationType,
			Lifetime:           Singleton, // All services are Singleton now
			ServiceKey:         key.Name,
		}
		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}

// register is a helper to register a singleton service.
func (s *serviceCollection) register(constructor interface{}, lifetime ServiceLifetime) error {
	if constructor == nil {
		return fmt.Errorf("constructor cannot be nil")
	}

	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return fmt.Errorf("constructor must be a function")
	}

	if ctorType.NumOut() == 0 || ctorType.NumOut() > 2 {
		return fmt.Errorf("constructor must return 1 or 2 values")
	}

	if ctorType.NumOut() == 2 {
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !ctorType.Out(1).Implements(errorType) {
			return fmt.Errorf("second return value must be error")
		}
	}

	returnType := ctorType.Out(0)

	reg := &internal.Registration{
		ServiceType:        returnType,
		ImplementationType: returnType,
		Lifetime:           internal.Singleton,
		Factory:            constructor,
	}

	return s.engine.Register(reg)
}

// registerKeyed is a helper to register a keyed singleton service.
func (s *serviceCollection) registerKeyed(constructor interface{}, lifetime ServiceLifetime, serviceKey string) error {
	if constructor == nil {
		return fmt.Errorf("constructor cannot be nil")
	}

	if serviceKey == "" {
		return fmt.Errorf("serviceKey cannot be empty")
	}

	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return fmt.Errorf("constructor must be a function")
	}

	if ctorType.NumOut() == 0 || ctorType.NumOut() > 2 {
		return fmt.Errorf("constructor must return 1 or 2 values")
	}

	if ctorType.NumOut() == 2 {
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !ctorType.Out(1).Implements(errorType) {
			return fmt.Errorf("second return value must be error")
		}
	}

	returnType := ctorType.Out(0)

	reg := &internal.Registration{
		ServiceType:        returnType,
		ImplementationType: returnType,
		Lifetime:           internal.Singleton,
		Factory:            constructor,
	}

	// Register with the service key in the engine
	return s.engine.RegisterKeyed(reg, serviceKey)
}
