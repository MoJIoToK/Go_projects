// Package logger provides structured JSON logging using log/slog.
// It configures a global default logger that writes to stdout with debug-level verbosity.
package logger

import (
	"log/slog"
	"os"
)

// SetupLogger initializes the global default logger with JSON output and debug level.
// After calling this function, slog.Default() returns a JSON-formatted logger.
func SetupLogger() {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(log)
}
