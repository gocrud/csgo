package di

import "reflect"

// ServiceDescriptor describes a service registration.
type ServiceDescriptor struct {
	// ServiceType is the type of the service (usually an interface).
	ServiceType reflect.Type

	// ImplementationType is the concrete type that implements the service.
	ImplementationType reflect.Type

	// Lifetime specifies when to create service instances.
	Lifetime ServiceLifetime

	// ServiceKey is used for keyed services (optional).
	ServiceKey string

	// Factory is the function that creates the service instance.
	Factory interface{}

	// Instance is a pre-created singleton instance (optional).
	Instance interface{}

	// Interfaces lists additional interface types this service implements.
	Interfaces []reflect.Type
}
