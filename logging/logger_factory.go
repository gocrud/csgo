package logging

// ILoggerFactory is used to create ILogger instances.
// Corresponds to .NET Microsoft.Extensions.Logging.ILoggerFactory.
type ILoggerFactory interface {
	// CreateLogger creates a logger with the specified category name.
	CreateLogger(category string) ILogger
}
