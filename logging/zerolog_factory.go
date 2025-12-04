package logging

import (
	"github.com/rs/zerolog"
)

// ZerologFactory is an implementation of ILoggerFactory using zerolog.
type ZerologFactory struct {
	logger   zerolog.Logger
	minLevel LogLevel
}

// NewZerologFactory creates a new ZerologFactory.
func NewZerologFactory(logger zerolog.Logger, minLevel LogLevel) *ZerologFactory {
	// Set global zerolog level
	zerolog.SetGlobalLevel(toZerologLevel(minLevel))

	return &ZerologFactory{
		logger:   logger,
		minLevel: minLevel,
	}
}

// CreateLogger creates a logger with the specified category name.
func (f *ZerologFactory) CreateLogger(category string) ILogger {
	return NewZerologLogger(f.logger, category, f.minLevel)
}
