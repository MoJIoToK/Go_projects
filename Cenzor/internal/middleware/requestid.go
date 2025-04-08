//Пакет middleware/requestID служит для генерации и использования уникального идентификатора
//запроса(Request ID) в HTTP-запросах. С помощью него происходит присвоение каждому запросу уникального
//идентификатора, который можно использовать для отслеживания запроса на протяжении его обработки сервером.
//Идентификатор запроса может использоваться для логирования, диагностики и мониторинга, когда нужно
//точно связать действия, происходящие в рамках одного запроса.

package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/sqids/sqids-go"
)

// ctxKey - тип ключа для ID запроса внутри контекста.
type ctxKey int

// RequestIDKey - ключ ID запроса внутри контекста.
const RequestIDKey ctxKey = 0

// RequestIDHeader - HTTP заголовок ID запроса.
const RequestIDHeader = "X-Request-Id"

// RequestID - функция проверяет наличие уникального ID запроса в заголовках и записывает
// значение в контекст. Если ID не найдено, то идёт генерация нового ID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			tm := time.Now()
			sec := uint64(tm.Unix())
			nano := uint64(tm.UnixNano())

			//Генерация короткого ID с помощью библиотеки sqids.
			s, _ := sqids.New(sqids.Options{MinLength: 10})
			id, err := s.Encode([]uint64{sec, nano})
			if err != nil {
				id = "unknown RequestID"
			}
			requestID = id

			r.Header.Set(RequestIDHeader, requestID)
		}

		//Сохранение Request ID в контекст запроса.
		//Это позволяет использовать Request ID в других частях программы через доступ к контексту.
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		//Вызов следующего обработчика
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID - функция возвращает ID запроса из контекста в виде строки.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}

	return ""
}
