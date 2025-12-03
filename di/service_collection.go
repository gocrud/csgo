package di

import (
	"fmt"
	"reflect"

	"github.com/gocrud/csgo/di/internal"
)

// IServiceCollection is a contract for a collection of service descriptors.
type IServiceCollection interface {
	// AddSingleton registers a singleton service.
	AddSingleton(constructor interface{}) IServiceCollection

	// AddTransient registers a transient service.
	AddTransient(constructor interface{}) IServiceCollection

	// TryAddSingleton attempts to add a singleton service if it doesn't exist.
	TryAddSingleton(constructor interface{}) IServiceCollection

	// TryAddTransient attempts to add a transient service if it doesn't exist.
	TryAddTransient(constructor interface{}) IServiceCollection

	// AddSingletonInstance registers a singleton instance.
	AddSingletonInstance(instance interface{}) IServiceCollection

	// AddHostedService registers a hosted service (background service).
	// Corresponds to .NET IServiceCollection.AddHostedService<T>().
	AddHostedService(constructor interface{}) IServiceCollection

	// AddKeyedSingleton registers a keyed singleton service.
	// Corresponds to .NET IServiceCollection.AddKeyedSingleton().
	AddKeyedSingleton(serviceKey string, constructor interface{}) IServiceCollection

	// AddKeyedTransient registers a keyed transient service.
	// Corresponds to .NET IServiceCollection.AddKeyedTransient().
	AddKeyedTransient(serviceKey string, constructor interface{}) IServiceCollection

	// BuildServiceProvider builds the service provider.
	BuildServiceProvider(options ...ServiceProviderOption) IServiceProvider

	// Count returns the number of registered services.
	Count() int

	// GetDescriptors returns all service descriptors.
	GetDescriptors() []ServiceDescriptor
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

// AddSingleton registers a singleton service.
func (s *serviceCollection) AddSingleton(constructor interface{}) IServiceCollection {
	if err := s.register(constructor, Singleton); err != nil {
		panic(fmt.Sprintf("failed to register singleton: %v", err))
	}
	return s
}

// AddTransient registers a transient service.
func (s *serviceCollection) AddTransient(constructor interface{}) IServiceCollection {
	if err := s.register(constructor, Transient); err != nil {
		panic(fmt.Sprintf("failed to register transient: %v", err))
	}
	return s
}

// TryAddSingleton attempts to add a singleton if it doesn't exist.
func (s *serviceCollection) TryAddSingleton(constructor interface{}) IServiceCollection {
	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return s
	}
	if ctorType.NumOut() == 0 {
		return s
	}

	returnType := ctorType.Out(0)
	if !s.engine.Contains(returnType, "") {
		s.AddSingleton(constructor)
	}
	return s
}

// TryAddTransient attempts to add a transient if it doesn't exist.
func (s *serviceCollection) TryAddTransient(constructor interface{}) IServiceCollection {
	ctorType := reflect.TypeOf(constructor)
	if ctorType.Kind() != reflect.Func {
		return s
	}
	if ctorType.NumOut() == 0 {
		return s
	}

	returnType := ctorType.Out(0)
	if !s.engine.Contains(returnType, "") {
		s.AddTransient(constructor)
	}
	return s
}

// AddSingletonInstance registers a singleton instance.
func (s *serviceCollection) AddSingletonInstance(instance interface{}) IServiceCollection {
	if instance == nil {
		panic("instance cannot be nil")
	}

	instanceType := reflect.TypeOf(instance)
	constructor := func() interface{} {
		return instance
	}

	reg := &internal.Registration{
		ServiceType:        instanceType,
		ImplementationType: instanceType,
		Lifetime:           internal.Singleton,
		Factory:            constructor,
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
	return s.AddSingleton(constructor)
}

// AddKeyedSingleton registers a keyed singleton service.
func (s *serviceCollection) AddKeyedSingleton(serviceKey string, constructor interface{}) IServiceCollection {
	if err := s.registerKeyed(constructor, Singleton, serviceKey); err != nil {
		panic(fmt.Sprintf("failed to register keyed singleton: %v", err))
	}
	return s
}

// AddKeyedTransient registers a keyed transient service.
func (s *serviceCollection) AddKeyedTransient(serviceKey string, constructor interface{}) IServiceCollection {
	if err := s.registerKeyed(constructor, Transient, serviceKey); err != nil {
		panic(fmt.Sprintf("failed to register keyed transient: %v", err))
	}
	return s
}

// BuildServiceProvider builds the service provider.
func (s *serviceCollection) BuildServiceProvider(options ...ServiceProviderOption) IServiceProvider {
	opts := &ServiceProviderOptions{}
	for _, opt := range options {
		opt(opts)
	}

	provider := &serviceProvider{
		engine: s.engine,
	}

	if err := s.engine.Compile(); err != nil {
		panic(fmt.Sprintf("failed to build service provider: %v", err))
	}

	return provider
}

// Count returns the number of registered services.
func (s *serviceCollection) Count() int {
	// TODO: implement actual count
	return 0
}

// GetDescriptors returns all service descriptors.
func (s *serviceCollection) GetDescriptors() []ServiceDescriptor {
	// TODO: implement actual descriptor retrieval
	return []ServiceDescriptor{}
}

// register is a helper to register a service with a specific lifetime.
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

	var internalLifetime internal.ServiceLifetime
	switch lifetime {
	case Singleton:
		internalLifetime = internal.Singleton
	case Transient:
		internalLifetime = internal.Transient
	default:
		internalLifetime = internal.Singleton
	}

	reg := &internal.Registration{
		ServiceType:        returnType,
		ImplementationType: returnType,
		Lifetime:           internalLifetime,
		Factory:            constructor,
	}

	return s.engine.Register(reg)
}

// registerKeyed is a helper to register a keyed service with a specific lifetime.
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

	var internalLifetime internal.ServiceLifetime
	switch lifetime {
	case Singleton:
		internalLifetime = internal.Singleton
	case Transient:
		internalLifetime = internal.Transient
	default:
		internalLifetime = internal.Singleton
	}

	reg := &internal.Registration{
		ServiceType:        returnType,
		ImplementationType: returnType,
		Lifetime:           internalLifetime,
		Factory:            constructor,
	}

	// Register with the service key in the engine
	return s.engine.RegisterKeyed(reg, serviceKey)
}

// ServiceProviderOptions contains options for building a service provider.
type ServiceProviderOptions struct {
	ValidateOnBuild bool
}

// ServiceProviderOption configures a ServiceProviderOptions.
type ServiceProviderOption func(*ServiceProviderOptions)

// WithValidateOnBuild validates all services can be resolved during build.
func WithValidateOnBuild(validate bool) ServiceProviderOption {
	return func(opts *ServiceProviderOptions) {
		opts.ValidateOnBuild = validate
	}
}
