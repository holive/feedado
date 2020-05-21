package feed

import "context"

type Repository interface {
	Create(ctx context.Context, feed *Feed) (*Feed, error)
	Update(ctx context.Context, feed *Feed) (*Feed, error)
	DeleteBySource(ctx context.Context, source string) error
	FindBySource(ctx context.Context, source string) (*Feed, error)
	FindAll(ctx context.Context, limit string, offset string) (SearchResult, error)
}
