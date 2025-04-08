//Пакет middleware/logger реализует промежуточный обработчик для веб-сервера, который логирует
//информацию о каждом запросе и ответе.

package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

// loggingResponseWriter - обертка над стандартным http.ResponseWriter для сохранения кода ответа.
// Т.к. с помощью оригинального интерфейса выполнить сохранение кода невозможно.
// Этот статус используется для записи в лог.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter - конструктор loggingResponseWriter.
// По умолчанию статус ответа 200.
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

// WriteHeader - метод записи кода в ответ.
// Статус ответа сохраняется в поле statusCode, а затем вызывается оригинальный метод WriteHeader.
func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

// Logger - функция, которая позволяет записывать логи запроса и ответа в логгер slog.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request log:",
			slog.String("host", r.Host),
			slog.String("uri", r.RequestURI),
			slog.String("protocol", r.Proto),
			slog.String("method", r.Method),
			slog.String("remote_address", r.RemoteAddr),
			slog.String("real_IP", r.Header.Get("X-Real-IP")),
			slog.String("request_id", GetRequestID(r.Context())),
		)

		//Создание нового loggingResponseWriter.
		lw := NewLoggingResponseWriter(w)

		//Вызов следующего обработчика
		next.ServeHTTP(lw, r)

		code := fmt.Sprintf("%d %s", lw.statusCode, http.StatusText(lw.statusCode))

		//Логирование статуса и тип контента
		slog.Info("Response log:",
			slog.String("status", code),
			slog.String("content-type", w.Header().Get("Content-Type")),
		)
	})
}
