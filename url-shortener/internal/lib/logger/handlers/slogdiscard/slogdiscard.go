package slogdiscard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct {
}

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (d DiscardHandler) Enabled(ctx context.Context, level slog.Level) bool {
	//Записи в журнал игнорируются, поэтому всегда false
	return true
}

func (d *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	//Возвращает тот же обработчик, т.к. нет атрибутов для сохранения
	return nil
}

func (d *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//Возвращает тот же обработчик, т.к. нет атрибутов для сохранения
	return d
}

func (d *DiscardHandler) WithGroup(name string) slog.Handler {
	//Возвращает тот же обработчик, т.к. нет атрибутов для сохранения
	return d
}
