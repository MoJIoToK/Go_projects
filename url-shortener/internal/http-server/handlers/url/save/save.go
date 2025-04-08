package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

// Структура запросов
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// Структура для ответов запросов
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLSaver
// TODO: перенести в конфиг
const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

// New - Конструктор хендлера для сохранения URL.
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to decode request"))

			return
		}
		log.Info("Request body decoded", slog.Any("request", req))

		//Создается валидатор который валидирует структуру. Проверяет корректность URL
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Failed to validate request", sl.Err(err))

			//Возвращает человекочитаемую ошибку
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("URL already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("URL already exists"))

			return
		}

		if err != nil {
			log.Error("Failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to add url"))

			return
		}

		log.Info("URL added", slog.Int64("id", id))

		ResponseOk(w, r, alias)

	}
}

func ResponseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
