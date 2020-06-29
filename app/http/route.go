package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/holive/feedado/app/http/handler"
	"github.com/rs/cors"
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
	r.Use(cors.AllowAll().Handler)

	r.Get("/health", handler.Health)

	r.Route("/feed", func(r chi.Router) {
		r.Post("/", handler.CreateFeed)
		r.Get("/", handler.GetAllFeeds)
		r.Get("/categories", handler.GetAllCategories)
		r.Get("/{source}", handler.GetFeed)
		r.Put("/", handler.UpdateFeed)
		r.Delete("/{source}", handler.DeleteFeed)
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/", handler.GetAllUsers)
		r.Get("/{email}", handler.GetUser)
		r.Put("/", handler.UpdateUser)
		r.Delete("/{email}", handler.DeleteUser)
	})

	r.Route("/rss", func(r chi.Router) {
		r.Get("/", handler.GetAllRSS)
		r.Get("/category/{category}", handler.GetAllRSSByCategory)
		r.Delete("/{url}", handler.DeleteRSS)
	})

	return r
}

func NewWorkerRouter(cfg *RouterConfig, handler *handler.WorkerHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(cfg.MiddlewareTimeout))
	r.Use(cors.Default().Handler)

	r.Get("/health", handler.Health)

	r.Route("/feedado-worker", func(r chi.Router) {
		r.Post("/feed", handler.ReindexFeeds)
		r.Post("/feed/category/{category}", handler.ReindexFeedsByCategory)
	})

	return r
}
