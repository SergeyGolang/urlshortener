package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"urlshortener/internal/storage"
	resp "urlshortener/lib/api/response"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("alias not found", slog.String("alias", alias))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}

		if err != nil {
			log.Error("failed to delete url", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}

		log.Info("deleted alias", slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
