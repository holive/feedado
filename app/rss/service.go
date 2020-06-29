package rss

import (
	"context"

	infraHTTP "github.com/holive/gopkg/net/http"
)

type Service struct {
	repo       Repository
	httpRunner infraHTTP.Runner
}

func (s *Service) FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *Service) Delete(ctx context.Context, url string) error {
	return s.repo.Delete(ctx, url)
}

func (s *Service) FindAllByCategory(ctx context.Context, limit string, offset string,
	category string) (*SearchResult, error) {

	return s.repo.FindAllByCategory(ctx, limit, offset, category)
}

func NewService(repository Repository, client infraHTTP.Runner) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
	}
}
