package hosting

import (
	"github.com/gocrud/csgo/configuration"
	"github.com/gocrud/csgo/di"
)

// IHostBuilder provides a mechanism for configuring and creating a host.
type IHostBuilder interface {
	ConfigureServices(configure func(services di.IServiceCollection)) IHostBuilder
	ConfigureAppConfiguration(configure func(config configuration.IConfigurationBuilder)) IHostBuilder
	ConfigureHostConfiguration(configure func(config configuration.IConfigurationBuilder)) IHostBuilder
	Build() IHost
}

// HostBuilder is the default implementation of IHostBuilder.
// Corresponds to .NET HostBuilder.
type HostBuilder struct {
	// Public properties (exposed like .NET)
	Services      di.IServiceCollection
	Configuration configuration.IConfiguration
	Environment   *Environment
}

// CreateDefaultBuilder creates a host builder with default configuration.
// Corresponds to .NET Host.CreateDefaultBuilder(args).
func CreateDefaultBuilder(args ...string) *HostBuilder {
	env := NewEnvironment()

	// Build configuration
	configBuilder := configuration.NewConfigurationBuilder().
		AddJsonFile("appsettings.json", true, true).
		AddJsonFile("appsettings."+env.Name()+".json", true, true).
		AddEnvironmentVariables("").
		AddCommandLine(args)

	config := configBuilder.Build()

	// Create service collection
	services := di.NewServiceCollection()

	// Register core services
	services.AddSingleton(func() configuration.IConfiguration { return config })
	services.AddSingleton(func() *Environment { return env })
	services.AddSingleton(func() IHostApplicationLifetime { return NewApplicationLifetime() })

	return &HostBuilder{
		Services:      services,
		Configuration: config,
		Environment:   env,
	}
}

// CreateEmptyBuilder creates an empty host builder without default configuration.
// Corresponds to .NET new HostBuilder().
func CreateEmptyBuilder() *HostBuilder {
	env := NewEnvironment()
	services := di.NewServiceCollection()

	return &HostBuilder{
		Services:    services,
		Environment: env,
	}
}

// ConfigureServices adds a delegate for configuring services.
// This method is optional - you can also directly access builder.Services.
func (b *HostBuilder) ConfigureServices(configure func(services di.IServiceCollection)) IHostBuilder {
	configure(b.Services)
	return b
}

// ConfigureAppConfiguration adds a delegate for configuring the application configuration.
func (b *HostBuilder) ConfigureAppConfiguration(configure func(config configuration.IConfigurationBuilder)) IHostBuilder {
	// For now, configuration is already built
	// In a full implementation, this would rebuild the configuration
	return b
}

// ConfigureHostConfiguration adds a delegate for configuring the host configuration.
func (b *HostBuilder) ConfigureHostConfiguration(configure func(config configuration.IConfigurationBuilder)) IHostBuilder {
	// For now, configuration is already built
	// In a full implementation, this would rebuild the configuration
	return b
}

// Build builds the host.
func (b *HostBuilder) Build() IHost {
	// Build service provider
	provider := b.Services.BuildServiceProvider()

	lifetime := NewApplicationLifetime()

	host := &Host{
		services:       provider,
		environment:    b.Environment,
		lifetime:       lifetime,
		hostedServices: make([]IHostedService, 0),
	}

	// Resolve hosted services
	hostedServices := b.resolveHostedServices(provider)
	host.hostedServices = hostedServices

	return host
}

// resolveHostedServices resolves all registered hosted services using new API
func (b *HostBuilder) resolveHostedServices(provider di.IServiceProvider) []IHostedService {
	var services []IHostedService

	// Use new pointer-filling API
	if err := provider.GetServices(&services); err != nil {
		// No hosted services registered or error occurred
		return []IHostedService{}
	}

	return services
}
