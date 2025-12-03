package di

// IServiceScopeFactory provides a mechanism to create service scopes.
// Corresponds to .NET IServiceScopeFactory.
type IServiceScopeFactory interface {
	// CreateScope creates a new IServiceScope.
	CreateScope() IServiceScope
}

// serviceScopeFactory is the default implementation of IServiceScopeFactory.
type serviceScopeFactory struct {
	provider IServiceProvider
}

// NewServiceScopeFactory creates a new service scope factory.
func NewServiceScopeFactory(provider IServiceProvider) IServiceScopeFactory {
	return &serviceScopeFactory{
		provider: provider,
	}
}

// CreateScope creates a new service scope.
func (f *serviceScopeFactory) CreateScope() IServiceScope {
	return f.provider.CreateScope()
}

// AddServiceScopeFactory registers the IServiceScopeFactory.
// This is automatically called by BuildServiceProvider.
func addServiceScopeFactory(services IServiceCollection, provider IServiceProvider) {
	factory := NewServiceScopeFactory(provider)
	services.AddSingletonInstance(factory)
}
