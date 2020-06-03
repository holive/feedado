package feed

import "context"

type WorkerRepository interface {
	Find(ctx context.Context, id string) (*Feed, error)
}
