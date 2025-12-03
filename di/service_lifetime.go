package di

// ServiceLifetime specifies the lifetime of a service in the dependency injection container.
type ServiceLifetime int

const (
	// Singleton specifies that a single instance of the service will be created.
	Singleton ServiceLifetime = iota

	// Transient specifies that a new instance of the service will be created every time it is requested.
	Transient
)

// String returns the string representation of the ServiceLifetime.
func (l ServiceLifetime) String() string {
	switch l {
	case Singleton:
		return "Singleton"
	case Transient:
		return "Transient"
	default:
		return "Unknown"
	}
}

