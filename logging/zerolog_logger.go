package logging

import (
	"fmt"

	"github.com/rs/zerolog"
)

// ZerologLogger is an implementation of ILogger using zerolog.
type ZerologLogger struct {
	logger   zerolog.Logger
	category string
	minLevel LogLevel
}

// NewZerologLogger creates a new ZerologLogger.
func NewZerologLogger(logger zerolog.Logger, category string, minLevel LogLevel) *ZerologLogger {
	return &ZerologLogger{
		logger:   logger.With().Str("category", category).Logger(),
		category: category,
		minLevel: minLevel,
	}
}

// Log writes a log entry at the specified log level.
func (l *ZerologLogger) Log(level LogLevel, message string, args ...interface{}) {
	if !l.IsEnabled(level) {
		return
	}

	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}

	switch level {
	case LogLevelTrace:
		l.logger.Trace().Msg(formattedMessage)
	case LogLevelDebug:
		l.logger.Debug().Msg(formattedMessage)
	case LogLevelInformation:
		l.logger.Info().Msg(formattedMessage)
	case LogLevelWarning:
		l.logger.Warn().Msg(formattedMessage)
	case LogLevelError:
		l.logger.Error().Msg(formattedMessage)
	case LogLevelCritical:
		l.logger.Fatal().Msg(formattedMessage)
	}
}

// LogTrace logs a trace message.
func (l *ZerologLogger) LogTrace(message string, args ...interface{}) {
	l.Log(LogLevelTrace, message, args...)
}

// LogDebug logs a debug message.
func (l *ZerologLogger) LogDebug(message string, args ...interface{}) {
	l.Log(LogLevelDebug, message, args...)
}

// LogInformation logs an informational message.
func (l *ZerologLogger) LogInformation(message string, args ...interface{}) {
	l.Log(LogLevelInformation, message, args...)
}

// LogWarning logs a warning message.
func (l *ZerologLogger) LogWarning(message string, args ...interface{}) {
	l.Log(LogLevelWarning, message, args...)
}

// LogError logs an error message.
func (l *ZerologLogger) LogError(err error, message string, args ...interface{}) {
	if !l.IsEnabled(LogLevelError) {
		return
	}

	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}

	if err != nil {
		l.logger.Error().Err(err).Msg(formattedMessage)
	} else {
		l.logger.Error().Msg(formattedMessage)
	}
}

// LogCritical logs a critical error message.
func (l *ZerologLogger) LogCritical(err error, message string, args ...interface{}) {
	if !l.IsEnabled(LogLevelCritical) {
		return
	}

	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}

	if err != nil {
		l.logger.Fatal().Err(err).Msg(formattedMessage)
	} else {
		l.logger.Fatal().Msg(formattedMessage)
	}
}

// IsEnabled checks if the given log level is enabled.
func (l *ZerologLogger) IsEnabled(level LogLevel) bool {
	return level >= l.minLevel && level != LogLevelNone
}

// toZerologLevel converts LogLevel to zerolog.Level.
func toZerologLevel(level LogLevel) zerolog.Level {
	switch level {
	case LogLevelTrace:
		return zerolog.TraceLevel
	case LogLevelDebug:
		return zerolog.DebugLevel
	case LogLevelInformation:
		return zerolog.InfoLevel
	case LogLevelWarning:
		return zerolog.WarnLevel
	case LogLevelError:
		return zerolog.ErrorLevel
	case LogLevelCritical:
		return zerolog.FatalLevel
	case LogLevelNone:
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}
