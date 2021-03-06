package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/holive/feedado/app/http/handler"

	"github.com/holive/feedado/app/feedado"
)

type Server struct {
	server *http.Server
}

func (s *Server) Start() error {
	fmt.Println("running...")
	return s.server.ListenAndServe()
}

type ServerConfig struct {
	Addr              string
	MaxHeaderBytes    int
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	Router            *RouterConfig
}

func NewServer(cfg *ServerConfig, services *feedado.Services) (*Server, error) {
	router := NewRouter(cfg.Router, &handler.Handler{
		Services: services,
	})

	return newServer(router, cfg), nil
}

func NewWorkerServer(cfg *ServerConfig, services *feedado.WorkerServices) (*Server, error) {
	router := NewWorkerRouter(cfg.Router, &handler.WorkerHandler{
		Services: services,
	})

	return newServer(router, cfg), nil
}

func newServer(router http.Handler, cfg *ServerConfig) *Server {
	return &Server{
		server: &http.Server{
			Handler:           router,
			Addr:              cfg.Addr,
			MaxHeaderBytes:    cfg.MaxHeaderBytes,
			IdleTimeout:       cfg.IdleTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
		},
	}
}
