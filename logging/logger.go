package logging

// ILogger represents a type used to perform logging.
// Corresponds to .NET Microsoft.Extensions.Logging.ILogger.
type ILogger interface {
	// Log writes a log entry at the specified log level.
	Log(level LogLevel, message string, args ...interface{})

	// LogTrace logs a trace message.
	LogTrace(message string, args ...interface{})

	// LogDebug logs a debug message.
	LogDebug(message string, args ...interface{})

	// LogInformation logs an informational message.
	LogInformation(message string, args ...interface{})

	// LogWarning logs a warning message.
	LogWarning(message string, args ...interface{})

	// LogError logs an error message.
	LogError(err error, message string, args ...interface{})

	// LogCritical logs a critical error message.
	LogCritical(err error, message string, args ...interface{})

	// IsEnabled checks if the given log level is enabled.
	IsEnabled(level LogLevel) bool
}
