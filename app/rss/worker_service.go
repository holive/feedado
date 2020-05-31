package rss

import (
	infraHTTP "github.com/holive/gopkg/net/http"
)

type WorkerService struct {
	repo       WorkerRepository
	httpRunner infraHTTP.Runner
}

func NewWorkerService(repository WorkerRepository, client infraHTTP.Runner) *WorkerService {
	return &WorkerService{
		repo:       repository,
		httpRunner: client,
	}
}
