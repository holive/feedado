package mongo

import (
	"context"
	"strconv"

	"github.com/holive/feedado/app/feed"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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

	if err = fr.collection.FindOne(ctx, bson.M{"_id": resp.InsertedID}).Decode(&newFeed); err != nil {
		return nil, errors.Wrap(err, "could not find the new feed")
	}

	return &newFeed, nil
}

func (fr *FeedRepository) Update(ctx context.Context, newFeed *feed.Feed) error {
	update, err := bson.Marshal(newFeed)
	if err != nil {
		return errors.Wrap(err, "could not marshal bson")
	}

	opts := options.Replace().SetUpsert(false)
	filter := bson.M{"source": bson.M{"$eq": newFeed.Source}}

	resp, err := fr.collection.ReplaceOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if resp.MatchedCount != 1 || resp.ModifiedCount != 1 {
		return errors.New("document not found or not updated")
	}

	return nil
}

func (fr *FeedRepository) Delete(ctx context.Context, source string) error {
	filter := bson.M{"source": bson.M{"$eq": source}}

	_, err := fr.collection.DeleteOne(ctx, filter)

	return err
}

func (fr *FeedRepository) FindBySource(ctx context.Context, source string) (*feed.Feed, error) {
	var f feed.Feed

	filter := bson.M{"source": bson.M{"$eq": source}}

	if err := fr.collection.FindOne(ctx, filter).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func (fr *FeedRepository) FindAll(ctx context.Context, limit string, offset string) (*feed.SearchResult, error) {
	intLimit, intOffset, err := fr.getLimitOffset(limit, offset)
	if err != nil {
		return &feed.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset)

	cur, err := fr.collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return &feed.SearchResult{}, err
	}

	total, err := fr.collection.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "could not count documents")
	}

	results, err := fr.resultFromCursor(ctx, cur)
	if err != nil {
		return &feed.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return &feed.SearchResult{
		Feeds: results,
		Result: feed.SearchResultResult{
			Offset: intOffset,
			Limit:  intLimit,
			Total:  total,
		},
	}, nil
}

func (fr *FeedRepository) FindAllCategories(ctx context.Context, limit string, offset string) (*feed.SearchResult, error) {
	intLimit, intOffset, err := fr.getLimitOffset(limit, offset)
	if err != nil {
		return &feed.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset).SetProjection(bson.M{"category": 1})

	cur, err := fr.collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return &feed.SearchResult{}, err
	}

	results, err := fr.resultFromCursor(ctx, cur)
	if err != nil {
		return &feed.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return &feed.SearchResult{
		Feeds:  results,
		Result: feed.SearchResultResult{},
	}, nil
}

func (fr *FeedRepository) getLimitOffset(limit string, offset string) (int64, int64, error) {
	if offset == "" {
		offset = "0"
	}

	if limit == "" {
		limit = "24"
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		return 0, 0, err
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return 0, 0, err
	}

	return int64(intLimit), int64(intOffset), nil
}

func (fr *FeedRepository) resultFromCursor(ctx context.Context, cur *mongo.Cursor) ([]feed.Feed, error) {
	var results []feed.Feed
	for cur.Next(ctx) {
		var elem feed.Feed
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	return results, nil
}

func NewFeedRepository(conn *Client) *FeedRepository {
	return &FeedRepository{
		collection: conn.db.Collection(FeedCollection),
	}
}
