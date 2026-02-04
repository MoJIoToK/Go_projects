package main

import (
	"context"
	"limits-app/internal/handlers"
	"limits-app/internal/logger"
	"limits-app/internal/normalize"
	"limits-app/internal/storage/sqlite"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logger.SetupLogger()
	db, err := sqlite.NewStorage("./db/limits.db")
	if err != nil {
		log.Fatal("DB init failed:", err)
	}

	normalizer := normalize.NewNormalizer()
	normalizeHandler := handlers.NewNormalizeHandler(normalizer, db, slog.Default())
	listHandler := handlers.NewListHandler(db, slog.Default())

	http.Handle("/limits/normalize", normalizeHandler)
	http.Handle("/limits", listHandler)

	srv := http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("Server started on :8080")
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal("Server failed:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Println("Got signal:", sig)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
