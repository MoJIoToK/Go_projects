package main

import (
	"Cenzor/internal/api"
	"Cenzor/internal/config"
	"Cenzor/internal/logger"
	"log/slog"
	"net/http"
)

// server - структура сервера, хранящая экземпляр БД и конфигурацию.
type server struct {
	api *api.API
	cfg *config.Config
}

func main() {

	//Инициализация логера.
	logger.SetupLogger()
	slog.Debug("Logger setup load successful")

	//Загрузка конфигураций из файла конфигурации.
	cfg := config.MustLoad("./config/config.yaml")
	slog.Debug("Load config file success")

	var s server

	s.cfg = cfg

	s.api = api.New(s.cfg)

	err := http.ListenAndServe(cfg.Address, s.api.Router())
	if err != nil {
		slog.Error("Failed to start server", err)
	}
}
