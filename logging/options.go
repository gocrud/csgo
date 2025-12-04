package logging

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// LoggingOptions represents logging configuration options.
type LoggingOptions struct {
	// MinLevel is the minimum log level to log.
	MinLevel LogLevel

	// Console configuration
	Console ConsoleOptions

	// File configuration
	File FileOptions
}

// ConsoleOptions represents console logging options.
type ConsoleOptions struct {
	// Enabled indicates whether console logging is enabled.
	Enabled bool

	// UseConsoleWriter uses human-readable console output (for development).
	UseConsoleWriter bool
}

// FileOptions represents file logging options.
type FileOptions struct {
	// Enabled indicates whether file logging is enabled.
	Enabled bool

	// Path is the path to the log file.
	Path string
}

// DefaultLoggingOptions returns the default logging options.
func DefaultLoggingOptions() *LoggingOptions {
	return &LoggingOptions{
		MinLevel: LogLevelInformation,
		Console: ConsoleOptions{
			Enabled:          true,
			UseConsoleWriter: false,
		},
		File: FileOptions{
			Enabled: false,
			Path:    "logs/app.log",
		},
	}
}

// CreateZerologLogger creates a zerolog.Logger based on the options.
func (opts *LoggingOptions) CreateZerologLogger() zerolog.Logger {
	var writers []io.Writer

	// Console output
	if opts.Console.Enabled {
		if opts.Console.UseConsoleWriter {
			// Human-readable console output (for development)
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			writers = append(writers, consoleWriter)
		} else {
			// JSON output (for production)
			writers = append(writers, os.Stdout)
		}
	}

	// File output
	if opts.File.Enabled && opts.File.Path != "" {
		// Create log file
		file, err := os.OpenFile(opts.File.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			writers = append(writers, file)
		}
		// Note: In production, consider using a proper file rotation library
	}

	// Create multi-writer
	var writer io.Writer
	if len(writers) == 0 {
		writer = os.Stdout // Fallback to stdout
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = zerolog.MultiLevelWriter(writers...)
	}

	// Create logger
	logger := zerolog.New(writer).With().Timestamp().Logger()

	// Set level
	logger = logger.Level(toZerologLevel(opts.MinLevel))

	return logger
}
