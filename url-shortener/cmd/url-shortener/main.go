package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "development"
	envProd  = "production"
)

func main() {

	cfg := config.MustLoad()

	//Инициализация логера
	log := setupLogger(cfg.Env)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	//Инициализация БД
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	//ИНициализация роутера
	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID)

	//TODO: run server
}

// Функция для конфигурации логера.
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
