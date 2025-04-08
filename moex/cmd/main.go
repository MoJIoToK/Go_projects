package main

import (
	"log/slog"
	"moex-app/internal/logger"
	"moex-app/internal/menu"
)

func main() {

	logger.SetupLogger()
	slog.Debug("Logger setup load successful")

	menu.StartMenu()

}
