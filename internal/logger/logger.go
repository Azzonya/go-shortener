// Package logger provides logging functionality for the URL shortener application.
package logger

import (
	"go.uber.org/zap"
)

// Log represents the global logger instance used throughout the application.
var Log = zap.NewNop()

// Initialize initializes the logger with the specified log level.
// It takes a log level string as input and sets up the logger accordingly.
// It returns an error if there is any issue initializing the logger.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
