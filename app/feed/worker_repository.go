package feed

import (
	"context"
)

type SchemaGetter interface {
	Find(ctx context.Context, id string) (*Feed, error)
}

type WorkerRepository interface {
	FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error)
	FindByCategory(ctx context.Context, limit string, offset string, category string) (*SearchResult, error)
}
