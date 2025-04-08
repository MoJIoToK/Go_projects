//Пакет api содержит обработчики api

package api

import (
	"GoNews/internal/logger"
	"GoNews/internal/middleware"
	"GoNews/internal/model"
	"GoNews/internal/storage"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API - программный интерфейс сервиса GoNews.
type API struct {
	r       *mux.Router
	storage storage.DB
}

// Response - основная структура ответа сервера.
type Response struct {
	Pagination Pagination   `json:"pagination"`
	News       []model.Post `json:"news"`
}

// Pagination - структура пагинации.
type Pagination struct {
	Total int `json:"total_pages"`
	Page  int `json:"current_page"`
	Limit int `json:"limit"`
}

// postPerPage - количество постов на страницу по умолчанию.
const postPerPage = 10

// New - конструктор API.
func New(storage storage.DB) *API {
	api := API{
		storage: storage,
		//Создание роутера
		r: mux.NewRouter(),
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
	api.r.HandleFunc("/news", api.PostsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/id/{id}", api.PostByID).Methods(http.MethodGet, http.MethodOptions)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./cmd/webapp"))))
	api.r.Use(middleware.Logger)
	api.r.Use(middleware.RequestID)
}

// PostsHandler - метод записывает в ResponseWrite ответ в формате JSON.
// В ответ включен объект пагинации и слайса постов из БД, которые соответствуют запросу.
func (api *API) PostsHandler(w http.ResponseWriter, r *http.Request) {
	const operation = "GoNews.API.PostsHandler"

	log := slog.Default().With(
		slog.String("op", operation),
		slog.String("request_id", middleware.GetRequestID(r.Context())),
	)

	log.Info("Request to receive posts")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()

	query := r.URL.Query().Get("s")
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	opt := &storage.Options{}
	if query != "" {
		opt.SearchQuery = query
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = postPerPage
	}

	num, err := api.storage.CountPosts(ctx, opt)
	if err != nil {
		log.Error("failed to count posts", logger.Err(err))
		http.Error(w, "failed to receive posts from DB", http.StatusInternalServerError)
		return
	}
	if num == 0 {
		log.Error("posts not found")
		http.Error(w, "posts not found", http.StatusNotFound)
		return
	}
	log.Debug("posts count successfully", slog.Int64("num", num))

	pgCount := int(num) / limit
	if int(num)%limit != 0 {
		pgCount++
	}
	if page > pgCount {
		page = 1
	}

	onPage := int(num) - (page-1)*limit
	if onPage > postPerPage {
		onPage = postPerPage
	}
	pg := Pagination{Total: pgCount, Page: page, Limit: onPage}

	opt.Count = limit
	opt.Offset = limit * (page - 1)

	posts, err := api.storage.GetPosts(ctx, opt)
	if err != nil {
		log.Error("failed to receive posts", logger.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debug("Posts received successfully", slog.Int("num", len(posts)))

	response := Response{
		Pagination: pg,
		News:       posts,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(response)
	if err != nil {
		log.Error("failed to encode posts", logger.Err(err))
		http.Error(w, "failed to encode posts", http.StatusInternalServerError)
		return
	}

	log.Info("request served successfully")

}

// PostByID - метод записывает в ResponseWriter один пост по переданному ID.
func (api *API) PostByID(w http.ResponseWriter, r *http.Request) {
	const operation = "GoNews.API.PostByID"

	log := slog.Default().With(
		slog.String("op", operation),
		slog.String("request_id", middleware.GetRequestID(r.Context())),
	)

	log.Info("request to receive post by ID")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		log.Error("empty post id")
		http.Error(w, "incorrect post id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	post, err := api.storage.PostByID(ctx, id)
	if err != nil {
		log.Error("failed to receive post by id", slog.String("id", id), logger.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, storage.ErrIncorrectId) {
			http.Error(w, "incorrect post id", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to receive post", http.StatusInternalServerError)
		return
	}
	log.Debug("post by ID received successfully", slog.String("id", id))

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err = enc.Encode(post)
	if err != nil {
		log.Error("failed to encode post", logger.Err(err))
		http.Error(w, "failed to encode post", http.StatusInternalServerError)
		return
	}

	log.Info("request served successfully", slog.String("id", id))
}
