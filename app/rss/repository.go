package rss

import "context"

type Repository interface {
	FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error)
	Delete(ctx context.Context, url string) error
	FindAllByCategory(ctx context.Context, limit string, offset string, category string) (*SearchResult, error)
}
