package logging

import (
	"log/slog"
	"os"
)

var loggingEnabled = false
var logger *slog.Logger

// SetupLogging sets up the logging.
// It allows logging to be enabled / disabled.
// It also allows logging to be configured as text or json.
func SetupLogging(enabled, isJson bool) {
	if enabled {
		if isJson {
			logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
		} else {
			logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
		}
	}
	loggingEnabled = enabled
}

// Info is a wrapper method for logging
func Info(msg string) {
	if loggingEnabled {
		logger.Info(msg)
	}
}

// Error is a wrapper method for logging
func Error(msg string) {
	if loggingEnabled {
		logger.Error(msg)
	}
}
