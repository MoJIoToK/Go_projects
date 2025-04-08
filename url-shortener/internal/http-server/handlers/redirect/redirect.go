package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

// New - конструктор хендлера для перенаправления по ссылке, храняшейся по её сокращению.
func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.new"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("URL not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("not found"))

			return

		}

		if err != nil {
			log.Error("Failed to get URL", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got URL", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)

	}

}
