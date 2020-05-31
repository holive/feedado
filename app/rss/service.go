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
	//if err := s.validateURL(feed.Source); err != nil {
	//	return &Feed{}, err
	//}
	//
	//feed.Source = strings.TrimSuffix(feed.Source, "/")
	//
	//alreadyExists, _ := s.repo.FindBySource(ctx, feed.Source)
	//if alreadyExists != nil {
	//	return &Feed{}, errors.New("source already exists")
	//}
	//
	//newFeed, err := s.repo.Create(ctx, feed)
	//if err != nil {
	//	return &Feed{}, errors.Wrap(err, "could not create a feed")
	//}

	return nil, nil
}

func NewService(repository Repository, client infraHTTP.Runner) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
	}
}
