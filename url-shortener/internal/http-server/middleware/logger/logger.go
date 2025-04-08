package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

// Опциональная штука.
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		//Создание копии логера
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		//Вывод сообщения об активации middleware при запуске приложения
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			//Выполняются до обработки запросов
			entry := log.With(
				//Составление лога запроса. Его метод, путь и т.д.
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				//Непосредственный вывод логов
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			//Передача управления следующему handler'у
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
