package mongo

import (
	"context"
	"strconv"

	"github.com/holive/feedado/app/rss"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RSSRepository struct {
	collection *mongo.Collection
}

func (rr *RSSRepository) FindAll(ctx context.Context, limit string, offset string) (*rss.SearchResult, error) {
	intLimit, intOffset, err := rr.getLimitOffset(limit, offset)
	if err != nil {
		return &rss.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset).SetSort(bson.D{{"timestamp", -1}})

	cur, err := rr.collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return &rss.SearchResult{}, err
	}

	total, err := rr.collection.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "could not count documents")
	}

	results, err := rr.resultFromCursor(ctx, cur)
	if err != nil {
		return &rss.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return &rss.SearchResult{
		Feeds: results,
		Result: rss.SearchResultResult{
			Offset: intOffset,
			Limit:  intLimit,
			Total:  total,
		},
	}, nil
}

func (rr *RSSRepository) Delete(ctx context.Context, url string) error {
	filter := bson.M{"url": bson.M{"$eq": url}}

	_, err := rr.collection.DeleteOne(ctx, filter)

	return err
}

func (rr *RSSRepository) FindAllByCategory(ctx context.Context, limit string, offset string, category string) (*rss.SearchResult, error) {
	intLimit, intOffset, err := rr.getLimitOffset(limit, offset)
	if err != nil {
		return &rss.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset).SetSort(bson.D{{"timestamp", -1}})

	filter := bson.M{"category": bson.M{"$eq": category}}

	cur, err := rr.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return &rss.SearchResult{}, err
	}

	total, err := rr.collection.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "could not count documents")
	}

	results, err := rr.resultFromCursor(ctx, cur)
	if err != nil {
		return &rss.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return &rss.SearchResult{
		Feeds: results,
		Result: rss.SearchResultResult{
			Offset: intOffset,
			Limit:  intLimit,
			Total:  total,
		},
	}, nil
}

func (rr *RSSRepository) getLimitOffset(limit string, offset string) (int64, int64, error) {
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

func (rr *RSSRepository) resultFromCursor(ctx context.Context, cur *mongo.Cursor) ([]rss.RSS, error) {
	var results []rss.RSS
	for cur.Next(ctx) {
		var elem rss.RSS
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

func NewRssRepository(conn *Client) *RSSRepository {
	return &RSSRepository{
		collection: conn.db.Collection(RssCollection),
	}
}
