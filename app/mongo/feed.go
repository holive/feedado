package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/holive/feedado/app/feed"
)

type FeedRepository struct {
	collection *mongo.Collection
}

func (fr *FeedRepository) Create(ctx context.Context, fd *feed.Feed) (*feed.Feed, error) {
	resp, err := fr.collection.InsertOne(ctx, fd)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a feed")
	}

	var newFeed feed.Feed

	err = fr.collection.FindOne(ctx, bson.M{"_id": resp.InsertedID}).Decode(&newFeed)
	if err != nil {
		return nil, errors.Wrap(err, "could not find the new feed")
	}

	return &newFeed, nil
}

func (fr *FeedRepository) Update(ctx context.Context, feed *feed.Feed) (*feed.Feed, error) {
	panic("implement me")
}

func (fr *FeedRepository) DeleteBySource(ctx context.Context, source string) error {
	panic("implement me")
}

func (fr *FeedRepository) FindBySource(ctx context.Context, source string) (*feed.Feed, error) {
	panic("implement me")
}

func (fr *FeedRepository) FindAll(ctx context.Context, limit string, offset string) (feed.SearchResult, error) {
	panic("implement me")
}

func NewFeedRepository(conn *Client) *FeedRepository {
	return &FeedRepository{
		collection: conn.db.Collection("feed"),
	}
}
