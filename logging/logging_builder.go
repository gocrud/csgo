package logging

import (
	"github.com/gocrud/csgo/configuration"
)

// LoggingBuilder allows configuring logging services.
type LoggingBuilder struct {
	options *LoggingOptions
	config  configuration.IConfiguration
}

// NewLoggingBuilder creates a new LoggingBuilder.
func NewLoggingBuilder(config configuration.IConfiguration) *LoggingBuilder {
	builder := &LoggingBuilder{
		options: DefaultLoggingOptions(),
		config:  config,
	}

	// Load from configuration if available
	builder.loadFromConfiguration()

	return builder
}

// SetMinimumLevel sets the minimum log level.
func (b *LoggingBuilder) SetMinimumLevel(level LogLevel) *LoggingBuilder {
	b.options.MinLevel = level
	return b
}

// AddConsole enables console logging.
func (b *LoggingBuilder) AddConsole(useConsoleWriter bool) *LoggingBuilder {
	b.options.Console.Enabled = true
	b.options.Console.UseConsoleWriter = useConsoleWriter
	return b
}

// AddFile enables file logging.
func (b *LoggingBuilder) AddFile(path string) *LoggingBuilder {
	b.options.File.Enabled = true
	b.options.File.Path = path
	return b
}

// GetOptions returns the configured logging options.
func (b *LoggingBuilder) GetOptions() *LoggingOptions {
	return b.options
}

// loadFromConfiguration loads logging configuration from IConfiguration.
// Supports reading from appsettings.json:
//   Logging:LogLevel:Default
//   Logging:Console:Enabled
//   Logging:File:Enabled
//   Logging:File:Path
func (b *LoggingBuilder) loadFromConfiguration() {
	if b.config == nil {
		return
	}

	// Read log level
	if levelStr := b.config.Get("Logging:LogLevel:Default"); levelStr != "" {
		b.options.MinLevel = ParseLogLevel(levelStr)
	}

	// Read console options
	if enabledStr := b.config.Get("Logging:Console:Enabled"); enabledStr != "" {
		b.options.Console.Enabled = enabledStr == "true"
	}

	// Read file options
	if enabledStr := b.config.Get("Logging:File:Enabled"); enabledStr != "" {
		b.options.File.Enabled = enabledStr == "true"
	}

	if path := b.config.Get("Logging:File:Path"); path != "" {
		b.options.File.Path = path
	}
}
