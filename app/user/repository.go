package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, email string) error
	Find(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error)
}
