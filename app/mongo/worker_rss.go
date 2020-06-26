package mongo

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/rss"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RssWorkerRepository struct {
	collection *mongo.Collection
}

func (rr *RssWorkerRepository) Create(ctx context.Context, feeds []*rss.RSS) error {
	opts := options.Replace().SetUpsert(true)

	for _, f := range feeds {
		update, err := bson.Marshal(f)
		if err != nil {
			return errors.Wrap(err, "could not marshal bson")
		}

		filter := bson.M{"source": bson.M{"$eq": f.URL}}

		_, err = rr.collection.ReplaceOne(ctx, filter, update, opts)
		if err != nil {
			return errors.Wrap(err, "could not insert / update rss")
		}
	}

	return nil
}

func NewRssWorkerRepository(conn *Client, cfg *config.Config, logger *zap.SugaredLogger) *RssWorkerRepository {
	collection := conn.db.Collection(RssCollection)

	keys := bsonx.Doc{{Key: "timestamp", Value: bsonx.Int32(int32(1))}}
	idx := mongo.IndexModel{Keys: keys, Options: &options.IndexOptions{ExpireAfterSeconds: &cfg.Mongo.RSSttl}}

	_, err := collection.Indexes().CreateOne(context.Background(), idx)
	if err != nil {
		logger.Error("Error occurred while creating index", err)
	}

	return &RssWorkerRepository{
		collection: collection,
	}
}
