//Пакет api содержит обработчики api

package api

import (
	"Cenzor/internal/config"
	"Cenzor/internal/logger"
	"Cenzor/internal/middleware"
	"Cenzor/internal/model"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// API - программный интерфейс сервиса Cenzor
type API struct {
	r   *mux.Router
	cfg *config.Config
}

// New - конструктор API.
func New(cfg *config.Config) *API {
	api := API{
		r:   mux.NewRouter(),
		cfg: cfg,
	}
	//Запись обработчиков в структуре API
	api.endpoints()
	return &api
}

// Router - метод возвращает маршрутизатор для использования в качестве
// аргумента HTTP-сервера.
func (api *API) Router() *mux.Router {
	return api.r
}

// endpoints - метод регистрирует методы API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.Use(middleware.RequestID)
	api.r.Use(middleware.Logger)
	api.r.HandleFunc("/", api.Cenzor).Methods(http.MethodPost, http.MethodOptions)
}

// Cenzor - метод для цензурирования комментариев. Возвращает код 200, если комментарий успешно прошел
// проверку либо код 400, если проверку не была пройдена.
func (api *API) Cenzor(w http.ResponseWriter, r *http.Request) {

	const operation = "goCenzor.server.Cenzor"

	log := slog.Default().With(
		slog.String("op", operation),
		slog.String("request_id", middleware.GetRequestID(r.Context())),
	)

	log.Info("Comment censor request")

	r.Body = http.MaxBytesReader(w, r.Body, api.cfg.MaxBodySize)
	defer r.Body.Close()

	var req model.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		
		// Проверяем, была ли ошибка вызвана превышением лимита на размер тела
		if strings.Contains(err.Error(), "http: request body too large") {
			// Логируем ошибку
			log.Error("Request body too large", logger.Err(err))
			// Возвращаем код 413 и сообщение об ошибке
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		log.Error("Request cannot be decoded", logger.Err(err))
		http.Error(w, "Request cannot be decoded", http.StatusBadRequest)
		return
	}

	if isProhibited(req.Content, api.cfg.CensorList) {
		log.Info("The comment contains prohibited words")
		http.Error(w, "The comment contains prohibited words", http.StatusBadRequest)
		return
	}
	log.Info("Comment is allowed")

	w.WriteHeader(http.StatusOK)
	log.Info("Request served successfully")
}

// isProhibited - функция проверяет текст контента комментария на содержание недопустимых выражений.
// Проверяется вхождение недопустимого выражения в текст.
func isProhibited(text string, words []string) bool {
	for _, word := range words {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}
