package main

import (
	"APIGateWay/internal/api"
	"APIGateWay/internal/config"
	"APIGateWay/internal/logger"
	"log"
	"log/slog"
	"net/http"
)

// server - структура сервера api.
type server struct {
	api *api.API
}

func main() {

	//Инициализация логера.
	logger.SetupLogger()
	slog.Debug("Logger setup load successful")

	//Загрузка конфигураций из файла конфигурации.
	cfg := config.MustLoad("./config/config.yaml")
	slog.Debug("Load config file success")

	var s server

	s.api = api.New(cfg)

	err := http.ListenAndServe(":80", s.api.Router())
	if err != nil {
		log.Fatal(err)
	}

}
