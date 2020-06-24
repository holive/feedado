package rss

import (
	"context"

	"github.com/holive/feedado/app/feed"
)

type Updater interface {
	Create(ctx context.Context, feeds []*RSS) error
}

type SchemaGetter interface {
	Find(ctx context.Context, id string) (*feed.Feed, error)
}

type WorkerRepository interface {
	FindAll(ctx context.Context, limit string, offset string) (*feed.SearchResult, error)
}
