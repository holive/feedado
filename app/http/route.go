package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/holive/feedado/app/http/handler"
)

type RouterConfig struct {
	MiddlewareTimeout time.Duration
}

func NewRouter(cfg *RouterConfig, handler *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(cfg.MiddlewareTimeout))

	r.Get("/health", handler.Health)

	r.Route("/feed", func(r chi.Router) {
		r.Post("/", handler.CreateFeed)
		r.Get("/", handler.GetAll)
		r.Get("/{source}", handler.GetFeed)
		r.Put("/", handler.UpdateFeed)
		r.Delete("/{source}", handler.DeleteFeed)
	})

	return r
}
