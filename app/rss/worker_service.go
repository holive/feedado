package rss

import (
	"context"
	"strconv"

	"github.com/holive/feedado/app/feed"
	"golang.org/x/sync/errgroup"

	"github.com/holive/feedado/app/gocloud"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type WorkerService struct {
	repo      WorkerRepository
	logger    *zap.SugaredLogger
	publisher *gocloud.RSSPublisher
}

func (ws *WorkerService) FindAll(ctx context.Context) error {
	bufferSize := 1 // TODO: get size from config
	g, ctx := errgroup.WithContext(ctx)
	feeds := make(chan feed.Feed, bufferSize)

	g.Go(func(c context.Context) func() error {
		return func() error {
			defer close(feeds)
			return ws.scrollFeeds(ctx, feeds)
		}
	}(ctx))

	for i := 0; i < bufferSize; i++ {
		g.Go(func(c context.Context) func() error {
			return func() error {
				return ws.work(ctx, feeds)
			}
		}(ctx))
	}

	return nil
}

func (ws *WorkerService) scrollFeeds(ctx context.Context, result chan feed.Feed) error {
	var limit = 1 // TODO: change value after test
	var offset = 0

	for {
		asdf := strconv.Itoa(limit)
		res, err := ws.repo.FindAll(ctx, asdf, strconv.Itoa(offset))
		if err != nil {
			return errors.Wrap(err, "could not get total schemas")
		}

		if len(res.Feeds) < limit {
			ws.logger.Info("exiting scroll")
			break
		}

		for _, f := range res.Feeds {
			select {
			case result <- f:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		offset = offset + limit
	}

	return nil
}

func (ws *WorkerService) work(ctx context.Context, feeds <-chan feed.Feed) error {
	for f := range feeds {
		err := ws.publisher.Publish(ctx, feed.FeedSQS{ID: f.Id.String()})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewWorkerService(repository WorkerRepository, logger *zap.SugaredLogger,
	publisher *gocloud.RSSPublisher) *WorkerService {

	return &WorkerService{
		repo:      repository,
		logger:    logger,
		publisher: publisher,
	}
}
