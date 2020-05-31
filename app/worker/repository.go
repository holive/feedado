package worker

import "context"

type Repository interface {
	Create(ctx context.Context, feed *Feed) (*Feed, error)
}
