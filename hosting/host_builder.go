package hosting

import (
	"strconv"
	"time"

	"github.com/gocrud/csgo/configuration"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/logging"
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
	Services             di.IServiceCollection
	Configuration        configuration.IConfigurationManager
	Environment          *Environment
	configurationActions []func(configuration.IConfigurationBuilder)
}

// CreateDefaultBuilder creates a host builder with default configuration.
// Corresponds to .NET Host.CreateDefaultBuilder(args).
func CreateDefaultBuilder(args ...string) *HostBuilder {
	env := NewEnvironment()

	// Create ConfigurationManager (allows dynamic configuration)
	configManager := configuration.NewConfigurationManager()

	// Add default configuration sources
	configManager.
		AddYamlFile("appsettings.yaml", true, true).
		AddYamlFile("appsettings."+env.Name()+".yaml", true, true).
		AddEnvironmentVariables("").
		AddCommandLine(args)

	// Create service collection
	services := di.NewServiceCollection()

	// Register core services
	services.Add(func() configuration.IConfiguration { return configManager })
	services.Add(func() configuration.IConfigurationManager { return configManager })
	services.Add(func() *Environment { return env })
	services.Add(func() IHostApplicationLifetime { return NewApplicationLifetime() })

	// Register logging services by default (like .NET)
	// This will read configuration from appsettings.json automatically
	addDefaultLogging(services, configManager, env)

	return &HostBuilder{
		Services:             services,
		Configuration:        configManager,
		Environment:          env,
		configurationActions: make([]func(configuration.IConfigurationBuilder), 0),
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
	// Support multiple calls, accumulate configuration actions
	b.configurationActions = append(b.configurationActions, configure)

	// Apply configuration action immediately to the Configuration
	configure(b.Configuration)

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
	provider := di.BuildServiceProvider(b.Services)

	// Get application lifetime
	lifetime := di.GetOr(provider, NewApplicationLifetime())

	// Resolve hosted services
	hostedServices := b.resolveHostedServices(provider)

	// Get shutdown timeout from configuration (default 30 seconds)
	shutdownTimeout := b.getShutdownTimeout()

	// Create host
	host := NewHostWithTimeout(provider, b.Environment, lifetime, hostedServices, shutdownTimeout)

	return host
}

// resolveHostedServices resolves all registered hosted services using new API
func (b *HostBuilder) resolveHostedServices(provider di.IServiceProvider) []IHostedService {
	// Use GetAll generic API
	services := di.GetAll[IHostedService](provider)
	return services
}

// getShutdownTimeout gets the shutdown timeout from configuration or returns default.
func (b *HostBuilder) getShutdownTimeout() time.Duration {
	if b.Configuration == nil {
		return 30 * time.Second
	}

	timeoutStr := b.Configuration.Get("server.shutdownTimeout")
	if timeoutStr == "" {
		return 30 * time.Second
	}

	// Parse timeout in seconds
	if seconds, err := strconv.Atoi(timeoutStr); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	return 30 * time.Second
}

// addDefaultLogging adds default logging services to the service collection.
// This is called automatically by CreateDefaultBuilder().
func addDefaultLogging(services di.IServiceCollection, config configuration.IConfiguration, env *Environment) {
	// Use zerolog-based logging with configuration support
	logging.AddLogging(services, config)

	// In development, we could customize further if needed
	// The configuration will be read from appsettings.json automatically
}
