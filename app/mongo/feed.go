package mongo

import (
	"context"

	"github.com/holive/feedado/app/feed"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (fr *FeedRepository) Update(ctx context.Context, newFeed *feed.Feed, feedID string) error {
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return errors.Wrap(err, "ObjectIDFromHex ERROR")
	}

	update, err := bson.Marshal(newFeed)
	if err != nil {
		return errors.Wrap(err, "could not marshal bson")
	}

	opts := options.Replace().SetUpsert(false)
	filter := bson.M{"_id": bson.M{"$eq": objID}}

	resp, err := fr.collection.ReplaceOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if resp.MatchedCount != 1 || resp.ModifiedCount != 1 {
		return errors.New("document not found or not updated")
	}

	return nil
}

func (fr *FeedRepository) Delete(ctx context.Context, feedID string) error {
	objID, err := primitive.ObjectIDFromHex(feedID)
	if err != nil {
		return errors.Wrap(err, "ObjectIDFromHex ERROR")
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}

	if _, err := fr.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (fr *FeedRepository) FindBySource(ctx context.Context, source string) (*feed.Feed, error) {
	var f feed.Feed

	filter := bson.M{"source": bson.M{"$eq": source}}

	if err := fr.collection.FindOne(ctx, filter).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func (fr *FeedRepository) FindAll(ctx context.Context, limit string, offset string) (feed.SearchResult, error) {
	panic("implement me")
}

func NewFeedRepository(conn *Client) *FeedRepository {
	return &FeedRepository{
		collection: conn.db.Collection("feed"),
	}
}
