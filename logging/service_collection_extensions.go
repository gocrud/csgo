package logging

import (
	"github.com/gocrud/csgo/configuration"
	"github.com/gocrud/csgo/di"
)

// AddLogging adds logging services to the service collection.
// This method is automatically called by HostBuilder.CreateDefaultBuilder().
// Corresponds to .NET services.AddLogging().
func AddLogging(services di.IServiceCollection, config configuration.IConfiguration) {
	// Create logging builder
	builder := NewLoggingBuilder(config)

	// Get options
	opts := builder.GetOptions()

	// Create zerolog logger and factory
	zerologLogger := opts.CreateZerologLogger()
	factory := NewZerologFactory(zerologLogger, opts.MinLevel)

	// Register ILoggerFactory as singleton
	services.AddSingleton(func() ILoggerFactory {
		return factory
	})
}

// AddZerolog adds zerolog-based logging to the service collection with custom options.
func AddZerolog(services di.IServiceCollection, configure func(*LoggingOptions)) {
	opts := DefaultLoggingOptions()
	if configure != nil {
		configure(opts)
	}

	// Create zerolog logger and factory
	zerologLogger := opts.CreateZerologLogger()
	factory := NewZerologFactory(zerologLogger, opts.MinLevel)

	// Register ILoggerFactory as singleton
	services.AddSingleton(func() ILoggerFactory {
		return factory
	})
}
