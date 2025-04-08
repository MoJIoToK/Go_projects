package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/delete"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
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
		log.Error("Failed to initialize database", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	//ИНициализация роутера
	router := chi.NewRouter()

	//middleware
	//Показывает id запроса. Чтобы в случае ошибки было легко найти то, что нужно.
	router.Use(middleware.RequestID)
	//Логирует все запросы
	router.Use(middleware.Logger)
	//Самодельный логгер запросов
	router.Use(mwLogger.New(log))
	//Восстановление работоспособности паники вызванной в каком-то handler
	router.Use(middleware.Recoverer)
	//
	router.Use(middleware.URLFormat)

	//Подключаем к роутеру handlerы
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("Starting server", slog.String("address", cfg.Address))

	//Создание сервера
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeOut,
	}

	//Запуск сервера
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
	}

	log.Error("Stopping server", slog.String("address", cfg.Address))
}

// Функция для конфигурации логера.
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
