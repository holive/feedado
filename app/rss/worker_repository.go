package rss

import (
	"context"
)

type Updater interface {
	Create(ctx context.Context, feeds []*RSS) error
}
