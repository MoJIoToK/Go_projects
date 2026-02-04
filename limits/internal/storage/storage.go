// Package storage defines the interface for normalized client limits.
// It abstracts database operations to allow swapping implementations (e.g., SQLite, PostgreSQL).
package storage

import (
	"limits-app/internal/models"
)

// Storage defines the contract for saving and loading normalized client limits.
type Storage interface {
	Save([]models.ClientLimits) error
	LoadAll() ([]models.ClientLimits, error)
}
