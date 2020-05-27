package feed

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	infraHTTP "github.com/holive/gopkg/net/http"
)

type Service struct {
	repo       Repository
	httpRunner infraHTTP.Runner
}

func (s *Service) Create(ctx context.Context, feed *Feed) (*Feed, error) {
	if err := s.validateURL(feed.Source); err != nil {
		return &Feed{}, err
	}

	feed.Source = strings.TrimSuffix(feed.Source, "/")

	alreadyExists, _ := s.repo.FindBySource(ctx, feed.Source)
	if alreadyExists != nil {
		return &Feed{}, errors.New("source already exists")
	}

	newFeed, err := s.repo.Create(ctx, feed)
	if err != nil {
		return &Feed{}, errors.Wrap(err, "could not create a feed")
	}

	return newFeed, nil
}

func (s *Service) Update(ctx context.Context, feed *Feed) error {
	return s.repo.Update(ctx, feed)
}

func (s *Service) Delete(ctx context.Context, source string) error {
	return s.repo.Delete(ctx, source)
}

func (s *Service) FindBySource(ctx context.Context, source string) (*Feed, error) {
	return s.repo.FindBySource(ctx, source)
}

func (s *Service) FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *Service) validateURL(source string) error {
	u, err := url.Parse(source)
	if (err == nil && u.Scheme != "" && u.Host != "") == false {
		return errors.New("invalid url")
	}

	if strings.HasPrefix(source, "https://") == false {
		return errors.New("source must have 'https://' as prefix") // TODO: is that right?
	}

	return nil
}

func NewService(repository Repository, client infraHTTP.Runner) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
	}
}
