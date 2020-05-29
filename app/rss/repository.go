package rss

import "context"

type Repository interface {
	Create(ctx context.Context, rss *RSS) (*RSS, error)
}
