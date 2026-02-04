// Package handlers provides HTTP request handlers for the limits normalization service.
// It implements REST endpoints for receiving and retrieving db.
package handlers

import (
	"io"
	"limits-app/internal/normalize"
	"limits-app/internal/parser"
	"limits-app/internal/storage"
	"log/slog"
	"net/http"
)

// NormalizeHandler is an HTTP handler that processes incoming limits db,
// normalizes it according to business rules, and saves the result to storage.
// It supports both file upload (multipart/form-db) and raw text input.
type NormalizeHandler struct {
	normalizer normalize.Normalizer
	storage    storage.Storage
	logger     *slog.Logger
}

// NewNormalizeHandler is a constructor for struct NormalizeHandler.
// It requires a configured normalize.Normalizer, storage.Storage, and slog.Logger.
func NewNormalizeHandler(
	normalizer normalize.Normalizer,
	storage storage.Storage,
	logger *slog.Logger) *NormalizeHandler {
	return &NormalizeHandler{
		normalizer: normalizer,
		storage:    storage,
		logger:     logger,
	}
}

// ServeHTTP handles POST /limits/normalize requests.
// It accepts either:
//   - A file upload with key "file" (multipart/form-db), or
//   - Raw text in the request body.
//
// The handler:
//  1. Parses DEPO records,
//  2. Normalizes them (removes zero-limit positions, adds fake ones if needed),
//  3. Validates presence of all 4 limits for non-zero positions,
//  4. Saves the result to the database.
//
// On success, it responds with HTTP 200 OK.
// On critical error (e.g., DB failure), it returns HTTP 500.
// Parsing and validation errors are logged but do not stop processing.
func (handler *NormalizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.NormalizeHandler.ServeHTTP"

	handler.logger.Info("Received normalization request")
	var reader io.Reader

	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("file")
	if err == nil && file != nil {
		defer file.Close()
		reader = file
	} else {
		reader = r.Body
	}

	depoLimits, parseErrs := parser.Parse(reader)
	if len(parseErrs) > 0 {
		for _, err := range parseErrs {
			handler.logger.Warn("Parse error",
				slog.String("op", op), slog.Any("error", err))
		}
	}

	input := normalize.Input{DepoLimits: depoLimits}
	output, err := handler.normalizer.Normalize(input)
	if err != nil {
		handler.logger.Warn("Normalize error",
			slog.String("op", op), slog.Any("error", err))
		http.Error(w, "Normalization failed", http.StatusInternalServerError)
		return
	}

	if len(output.Errors) > 0 {
		for _, err := range output.Errors {
			handler.logger.Warn("Validation error",
				slog.String("op", op), slog.Any("error", err))
		}
	}

	if err := handler.storage.Save(output.ClientLimits); err != nil {
		handler.logger.Warn("Failed to save to DB",
			slog.String("op", op), slog.Any("error", err))
		http.Error(w, "Failed to save to DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
