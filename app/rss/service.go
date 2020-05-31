package rss

import (
	"context"

	infraHTTP "github.com/holive/gopkg/net/http"
)

type Service struct {
	repo       Repository
	httpRunner infraHTTP.Runner
}

func (s *Service) Create(ctx context.Context, feed *RSS) (*RSS, error) {
	panic("create me")
}

func NewService(repository Repository, client infraHTTP.Runner) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
	}
}
