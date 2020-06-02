package main

import (
	"context"
	"fmt"
	"os"

	"github.com/holive/feedado/app/http"

	"github.com/holive/feedado/app/feedado"
	"github.com/pkg/errors"
)

func main() {
	worker, err := feedado.NewWorker()
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not run Feedado worker").Error())
		os.Exit(1)
	}

	err = worker.Worker.Start(context.Background())

	server, err := http.NewWorkerServer(&http.ServerConfig{
		Addr:              worker.Cfg.HTTPServer.Addr,
		MaxHeaderBytes:    worker.Cfg.HTTPServer.MaxHeaderBytes,
		IdleTimeout:       worker.Cfg.HTTPServer.IdleTimeout,
		ReadHeaderTimeout: worker.Cfg.HTTPServer.ReadHeaderTimeout,
		ReadTimeout:       worker.Cfg.HTTPServer.ReadTimeout,
		WriteTimeout:      worker.Cfg.HTTPServer.WriteTimeout,
		Router:            &http.RouterConfig{MiddlewareTimeout: worker.Cfg.HTTPServer.Router.MiddlewareTimeout},
	}, worker.Services)
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not run Feedado Worker").Error())
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		fmt.Println(errors.Wrap(err, "could not run Feedado Worker").Error())
		os.Exit(1)
	}
}
