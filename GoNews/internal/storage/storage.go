//Пакет содержит интерфейс для работы с любой реализацией базы данных, удовлетворяющей этому интерфейсу.

package storage

import (
	"GoNews/internal/model"
	"context"
	"errors"
)

// Ошибки при работе с БД.
var (
	ErrIncorrectId = errors.New("incorrect id")
	ErrNotFound    = errors.New("post not found")
)

// Options - структура запроса для текстового поиска в БД по заголовкам новостей.
type Options struct {
	// SearchQuery - запрос для текстового поиска.
	SearchQuery string

	// Count - максимальное число возвращаемых постов.
	Count int

	// Offset - число постов на сдвиг в пагинации.
	Offset int
}

// Interface - интерфейс хранилища постов из RSS лент.
type DB interface {
	AddPost(ctx context.Context, post []model.Post) (int, error)
	GetPosts(ctx context.Context, op ...*Options) ([]model.Post, error)
	PostByID(ctx context.Context, id string) (model.Post, error)
	CountPosts(ctx context.Context, q ...*Options) (int64, error)
	Close() error
}
