package mongo

import (
	"context"

	"github.com/holive/feedado/app/rss"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type RssWorkerRepository struct {
	collection *mongo.Collection
}

func (rr *RssWorkerRepository) Create(ctx context.Context, feeds []rss.RSS) error {

	_, err := rr.collection.InsertMany(ctx, feeds)
	if err != nil {
		return errors.Wrap(err, "could not create rss feeds")
	}

	return nil
}

func NewRssWorkerRepository(conn *Client) *RssWorkerRepository {
	return &RssWorkerRepository{
		collection: conn.db.Collection("rss"),
	}
}
