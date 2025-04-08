package api

import (
	"Comments/internal/config"
	"Comments/internal/logger"
	"Comments/internal/middleware"
	"Comments/internal/model"
	"Comments/internal/storage"
	"Comments/internal/tree"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strings"
)

// API - программный интерфейс сервиса GoNews
type API struct {
	r       *mux.Router
	storage storage.DB
	cfg     *config.Config
}

// New - конструктор API.
func New(storage storage.DB, cfg *config.Config) *API {
	api := API{
		storage: storage,
		//Создание роутера
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
	api.r.Use(middleware.Logger)
	api.r.Use(middleware.RequestID)
	api.r.HandleFunc("/news/comments/new", api.Add).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/news/id/{id}", api.Comments).Methods(http.MethodGet, http.MethodOptions)
}

// Add - метод позволяет добавить комментарий в БД. Тело запроса проверяется на наличие в поле content
// содержимого, если данное поле пустое, то в качестве ответа на запрос отправляется ошибка.
func (api *API) Add(w http.ResponseWriter, r *http.Request) {
	const operation = "goComments.server.AddComment"

	log := slog.Default().With(
		slog.String("op", operation),
		slog.String("request_id", middleware.GetRequestID(r.Context())),
	)

	log.Info("Add comment request")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	ct := r.Header.Get("Content-Type")
	media := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	if media != "application/json" {
		log.Error("Content-Type header is not application/json")
		http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, api.cfg.MaxBodySize)
	defer r.Body.Close()

	var comment model.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		log.Error("Request cannot be decoded", logger.Err(err))
		http.Error(w, "Request cannot be decoded", http.StatusBadRequest)
		return
	}
	log.Debug("request body decoded")

	if comment.Content == "" {
		log.Error("Comment has empty content field")
		http.Error(w, "Empty comment", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := api.storage.AddComment(ctx, comment)
	if err != nil {
		log.Error("Cannot add comment to DB", logger.Err(err))
		if errors.Is(err, storage.ErrIncorrectParentID) || errors.Is(err, storage.ErrIncorrectPostID) {
			http.Error(w, "Incorrect data", http.StatusBadRequest)
			return
		}
		http.Error(w, "Comment cannot be added", http.StatusInternalServerError)
		return
	}
	log.Debug("Comment added to DB successfully", slog.String("id", id))

	w.WriteHeader(http.StatusCreated)
	log.Info("Request served successfully")
}

// Comments - метод позволяет вывести все комментарии к определенной новости по её ID в формате ObjectID.
func (api *API) Comments(w http.ResponseWriter, r *http.Request) {
	const operation = "goComments.server.Comments"

	log := slog.Default().With(
		slog.String("op", operation),
		slog.String("request_id", middleware.GetRequestID(r.Context())),
	)

	log.Info("Receive comments request")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		log.Error("Empty news id")
		http.Error(w, "Empty news id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	comments, err := api.storage.Comments(ctx, id)
	if err != nil {
		log.Error("Cannot receive comments", logger.Err(err))
		if errors.Is(err, storage.ErrNoComments) {
			http.Error(w, "News id not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, storage.ErrIncorrectPostID) {
			http.Error(w, "Incorrect post id", http.StatusBadRequest)
			return
		}
		http.Error(w, "Cannot receive comments", http.StatusInternalServerError)
		return
	}

	log.Debug("Comments received successfully")

	root, err := tree.Build(comments)
	if err != nil {
		log.Error("Comments tree cannot be build", logger.Err(err))
		http.Error(w, "Comments tree cannot be build", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(root.Comments)
	if err != nil {
		log.Error("Comments cannot be encoded", logger.Err(err))
		http.Error(w, "Comments cannot be encoded", http.StatusInternalServerError)
		return
	}

	log.Info("Request served successfully")

}
