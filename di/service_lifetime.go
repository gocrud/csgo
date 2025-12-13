package di

// ServiceLifetime specifies the lifetime of a service in the dependency injection container.
// In this simplified DI container, only Singleton lifetime is supported.
type ServiceLifetime int

const (
	// Singleton specifies that a single instance of the service will be created
	// and shared across the entire application lifetime.
	Singleton ServiceLifetime = iota
)

// String returns the string representation of the ServiceLifetime.
func (l ServiceLifetime) String() string {
	return "Singleton"
}
