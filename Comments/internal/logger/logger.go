//Пакет logger для инициализации логгера и обертки для ошибок.

package logger

import (
	"log/slog"
	"os"
)

// SetupLogger - инициализирует логгер из пакета slog, где вывод осуществляется в формате JSON.
// Данный логгер устанавливает по умолчанию, что позволяет не передавать его другим объектам.
func SetupLogger() {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(log)
}

// Err - обертка для ошибки, представляет ее как атрибут логгера.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
