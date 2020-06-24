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
	bufferSize := 5
	g, ctx := errgroup.WithContext(ctx)
	feeds := make(chan feed.Feed, bufferSize)

	g.Go(func(c context.Context) func() error {
		return func() error {
			defer close(feeds)
			if err := ws.scrollFeeds(ctx, feeds); err != nil {
				ws.logger.Error(err)
			}
			return nil
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
	var limit = 30
	var offset = 0

	for {
		asdf := strconv.Itoa(limit)
		res, err := ws.repo.FindAll(ctx, asdf, strconv.Itoa(offset))
		if err != nil {
			return errors.Wrap(err, "could not get total schemas")
		}

		for _, f := range res.Feeds {
			select {
			case result <- f:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if len(res.Feeds) < limit {
			ws.logger.Info("finishing scroll")
			break
		}

		offset = offset + limit
	}

	return nil
}

func (ws *WorkerService) work(ctx context.Context, feeds <-chan feed.Feed) error {
	for f := range feeds {
		err := ws.publisher.Publish(ctx, feed.FeedSQS{ID: f.Id.Hex()})
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
