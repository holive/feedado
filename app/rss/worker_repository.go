package rss

import (
	"context"

	"github.com/holive/feedado/app/feed"
)

type WorkerRepository interface {
}

type SchemaGetter interface {
	Find(ctx context.Context, id string) (*feed.Feed, error)
}
