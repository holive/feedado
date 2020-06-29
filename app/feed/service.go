package feed

import (
	"context"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Service struct {
	repo Repository
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

func (s *Service) FindAllCategories(ctx context.Context) ([]string, error) {
	res, err := s.repo.FindAllCategories(ctx, "1000", "0")
	if err != nil {
		return nil, errors.Wrap(err, "could not get all categories")
	}

	m := make(map[string]interface{})
	for _, doc := range res.Feeds {
		m[doc.Category] = nil
	}

	var cats []string
	for key, _ := range m {
		cats = append(cats, key)
	}

	return cats, nil
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

func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}
