//Пакет содержит интерфейс для работы с любой реализацией базы данных, удовлетворяющей этому интерфейсу.

package storage

import (
	"Comments/internal/model"
	"context"
	"errors"
)

// Ошибки при работе с БД.
var (
	ErrNoComments        = errors.New("No comments on provided post id")
	ErrParentNotFound    = errors.New("Parent comment not found")
	ErrIncorrectParentID = errors.New("Incorrect parent id")
	ErrIncorrectPostID   = errors.New("Incorrect post id")
)

// DB Interface - интерфейс хранилища комментариев.
type DB interface {
	AddComment(ctx context.Context, comment model.Comment) (string, error)
	Comments(ctx context.Context, news string) ([]model.Comment, error)
}
