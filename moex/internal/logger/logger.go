package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func SetupLogger() {

	logFile, err := os.OpenFile("internal/logger/logger.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		slog.New(slog.NewTextHandler(os.Stdout, nil)).Error("error opening log file: %v", err)
		os.Exit(1)
	}

	log := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)
}
