package user

import (
	"context"

	"github.com/pkg/errors"
)

type Service struct {
	repo Repository
}

func (s *Service) Create(ctx context.Context, user *User) (*User, error) {
	alreadyExists, _ := s.repo.Find(ctx, user.Email)
	if alreadyExists != nil {
		return &User{}, errors.New("email already exists")
	}

	user, err := s.repo.Create(ctx, user)
	if err != nil {
		return &User{}, errors.Wrap(err, "could not create a feed")
	}

	return user, nil
}

func (s *Service) Update(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, email string) error {
	return s.repo.Delete(ctx, email)
}

func (s *Service) Find(ctx context.Context, email string) (*User, error) {
	return s.repo.Find(ctx, email)
}

func (s *Service) FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}
