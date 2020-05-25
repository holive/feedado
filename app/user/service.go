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

func (s *Service) Update(ctx context.Context, user *User, id string) error {
	return s.repo.Update(ctx, user, id)
}
