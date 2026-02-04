// Package handlers provides HTTP request handlers for the limits normalization service.
// It implements REST endpoints for receiving and retrieving db.
package handlers

import (
	"encoding/json"
	"limits-app/internal/storage"
	"log/slog"
	"net/http"
)

// ListHandler is HTTP handler that returns all normalized client limits from storage.
// It serves GET /limits requests and return db in JSON format.
type ListHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

// NewListHandler is a constructor for struct ListHandler.
// It requires a configured storage.Storage and a slog.Logger
func NewListHandler(storage storage.Storage, logger *slog.Logger) *ListHandler {
	return &ListHandler{
		storage: storage,
		logger:  logger,
	}
}

// ServeHTTP handles GET /limits requests.
// It retrieves all normalized client limits from storage and returns them as JSON.
// On success, it responds with HTTP 200 and JSON array.
// On error, it logs the failure and responds with HTTP 500.
func (handler *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.ListHandler"

	handler.logger.Info("Received request to list all limits")

	limits, err := handler.storage.LoadAll()
	if err != nil {
		handler.logger.Error("Failed to load limits from DB",
			slog.String("op", op), slog.Any("error", err))
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(limits); err != nil {
		handler.logger.Error("Failed to encode response",
			slog.String("op", op), slog.Any("error", err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	handler.logger.Info("Successfully returned all limits", slog.Int("count", len(limits)))
}
