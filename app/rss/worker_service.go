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
	repo      feed.WorkerRepository
	logger    *zap.SugaredLogger
	publisher *gocloud.RSSPublisher
}

type ScrollRepoCallback func(ctx context.Context, limit string, offset string) (*feed.SearchResult, error)

// this recursive function is necessary when one implement a reindex by source, for example.
// FindAllFeeds calls feedScrollPub passing a repo function that returns all feeds
func (ws *WorkerService) FindAllFeeds(ctx context.Context) error {
	return ws.feedScrollPub(ctx, func(ctx context.Context, limit string, offset string) (*feed.SearchResult, error) {
		return ws.repo.FindAll(ctx, limit, offset)
	})
}

// FindFeedByCategory calls feedScrollPub passing a repo function that returns a feeds with a specific source
func (ws *WorkerService) FindFeedByCategory(ctx context.Context, category string) error {
	return ws.feedScrollPub(ctx, func(ctx context.Context, limit string, offset string) (*feed.SearchResult, error) {
		return ws.repo.FindByCategory(ctx, limit, offset, category)
	})
}

// feedScrollPub triggers the scroll and publish jobs
func (ws *WorkerService) feedScrollPub(ctx context.Context, callback ScrollRepoCallback) error {
	bufferSize := 5
	g, ctx := errgroup.WithContext(ctx)
	feeds := make(chan feed.Feed, bufferSize)

	g.Go(func(c context.Context) func() error {
		return func() error {
			defer close(feeds)
			if err := ws.scrollFeeds(ctx, feeds, callback); err != nil {
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

func (ws *WorkerService) scrollFeeds(ctx context.Context, result chan feed.Feed, callback ScrollRepoCallback) error {
	var limit = 30
	var offset = 0

	for {
		res, err := callback(ctx, strconv.Itoa(limit), strconv.Itoa(offset))
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
			ws.logger.Debug("finishing scroll. length: ", len(res.Feeds))
			break
		}

		offset = offset + limit
	}

	return nil
}

func (ws *WorkerService) work(ctx context.Context, feeds <-chan feed.Feed) error {
	for f := range feeds {
		ws.logger.Debug("publishing ", f.Id)

		err := ws.publisher.Publish(ctx, feed.SQS{ID: f.Id.Hex()})
		if err != nil {
			return err
		}

		ws.logger.Debug("finish publishing ", f.Id)
	}
	return nil
}

func NewWorkerService(repository feed.WorkerRepository, logger *zap.SugaredLogger,
	publisher *gocloud.RSSPublisher) *WorkerService {

	return &WorkerService{
		repo:      repository,
		logger:    logger,
		publisher: publisher,
	}
}
