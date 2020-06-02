package rss

import (
	"context"

	infraHTTP "github.com/holive/gopkg/net/http"
)

type WorkerService struct {
	repo       Updater
	httpRunner infraHTTP.Runner
}

func (w WorkerService) Create(ctx context.Context, feeds []*RSS) error {
	panic("implement me")
}

func NewWorkerService(repository Updater, client infraHTTP.Runner) *WorkerService {
	return &WorkerService{
		repo:       repository,
		httpRunner: client,
	}
}
