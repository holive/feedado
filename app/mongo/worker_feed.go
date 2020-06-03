package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/holive/feedado/app/feed"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type FeedWorkerRepository struct {
	collection *mongo.Collection
}

func (fr *FeedWorkerRepository) Find(ctx context.Context, id string) (*feed.Feed, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get object id")
	}

	var f feed.Feed
	if err := fr.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func NewFeedWorkerRepository(conn *Client) *FeedWorkerRepository {
	return &FeedWorkerRepository{
		collection: conn.db.Collection(FeedCollection),
	}
}
