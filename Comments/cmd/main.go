package main

import (
	"Comments/internal/api"
	"Comments/internal/config"
	"Comments/internal/logger"
	"Comments/internal/storage"
	"Comments/internal/storage/mongo"
	"log"
	"log/slog"
	"net/http"
)

// server - структура сервера, хранящая экземпляр БД и api.
type server struct {
	db  storage.DB
	api *api.API
	cfg *config.Config
}

func main() {

	//Инициализация логера.
	logger.SetupLogger()
	slog.Debug("Logger setup load successful")

	var srv server

	//Загрузка конфигураций из файла конфигурации.
	cfg := config.MustLoad("./config/config.yaml")
	slog.Debug("Load config file success")

	//Инициализация зависимостей
	db, err := mongo.New(cfg.StoragePath)
	if err != nil {
		log.Fatal(err)
	}

	srv.db = db
	srv.api = api.New(srv.db, cfg)

	//Запуск веб-сервера с API и приложением
	err = http.ListenAndServe(cfg.Address, srv.api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
