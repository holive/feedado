package feed

import (
	"github.com/pkg/errors"

	infraHTTP "github.com/holive/gopkg/net/http"
)

type Service struct {
	repo       Repository
	httpRunner infraHTTP.Runner
}

func (s *Service) Create(feed *Feed) (*Feed, error) {
	// validar sources duplicados

	feed, err := s.repo.Create(feed)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a feed")
	}

	return feed, nil
}

func (s *Service) Update(feed *Feed) (*Feed, error) {
	// TODO: validar campos?

	feed, err := s.repo.Update(feed)
	if err != nil {
		return nil, errors.Wrap(err, "could not update the feed")
	}

	return feed, nil
}

func (s *Service) DeleteBySource(source string) error {
	return s.repo.DeleteBySource(source)
}

func (s *Service) FindBySource(source string) (*Feed, error) {
	return s.repo.FindBySource(source)
}

func (s *Service) FindAll(limit string, offset string) (SearchResult, error) {
	return s.repo.FindAll(limit, offset)
}

func NewService(repository Repository, client infraHTTP.Runner) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
	}
}
