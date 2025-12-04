package logging

// LogLevel represents the severity level of a log entry.
// Corresponds to .NET Microsoft.Extensions.Logging.LogLevel.
type LogLevel int

const (
	// LogLevelTrace logs that contain the most detailed messages.
	// These messages may contain sensitive application data and should not be enabled in production.
	LogLevelTrace LogLevel = iota

	// LogLevelDebug logs that are used for interactive investigation during development.
	LogLevelDebug

	// LogLevelInformation logs that track the general flow of the application.
	LogLevelInformation

	// LogLevelWarning logs that highlight an abnormal or unexpected event in the application flow.
	LogLevelWarning

	// LogLevelError logs that highlight when the current flow of execution is stopped due to a failure.
	LogLevelError

	// LogLevelCritical logs that describe an unrecoverable application or system crash.
	LogLevelCritical

	// LogLevelNone not used for writing log messages. Specifies that no messages should be written.
	LogLevelNone
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case LogLevelTrace:
		return "Trace"
	case LogLevelDebug:
		return "Debug"
	case LogLevelInformation:
		return "Information"
	case LogLevelWarning:
		return "Warning"
	case LogLevelError:
		return "Error"
	case LogLevelCritical:
		return "Critical"
	case LogLevelNone:
		return "None"
	default:
		return "Unknown"
	}
}

// ParseLogLevel parses a string into a LogLevel.
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "Trace":
		return LogLevelTrace
	case "Debug":
		return LogLevelDebug
	case "Information", "Info":
		return LogLevelInformation
	case "Warning", "Warn":
		return LogLevelWarning
	case "Error":
		return LogLevelError
	case "Critical":
		return LogLevelCritical
	case "None":
		return LogLevelNone
	default:
		return LogLevelInformation
	}
}
