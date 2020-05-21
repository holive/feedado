package feed

import (
	"context"

	"github.com/pkg/errors"

	infraHTTP "github.com/holive/gopkg/net/http"
)

type Service struct {
	repo       Repository
	httpRunner infraHTTP.Runner
}

func (s *Service) Create(ctx context.Context, feed *Feed) (*Feed, error) {
	// validar sources duplicados

	feed, err := s.repo.Create(ctx, feed)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a feed")
	}

	return feed, nil
}

func (s *Service) Update(ctx context.Context, feed *Feed) (*Feed, error) {
	// TODO: validar campos?

	feed, err := s.repo.Update(ctx, feed)
	if err != nil {
		return nil, errors.Wrap(err, "could not update the feed")
	}

	return feed, nil
}

func (s *Service) DeleteBySource(ctx context.Context, source string) error {
	return s.repo.DeleteBySource(ctx, source)
}

func (s *Service) FindBySource(ctx context.Context, source string) (*Feed, error) {
	return s.repo.FindBySource(ctx, source)
}

func (s *Service) FindAll(ctx context.Context, limit string, offset string) (SearchResult, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func NewService(repository Repository, client infraHTTP.Runner) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
	}
}
