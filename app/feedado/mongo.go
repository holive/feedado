package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/mongo"
)

func initMongoClient(cfg *config.Config) (*mongo.Client, error) {
	return mongo.New(&mongo.ClientConfig{
		URI:      cfg.Mongo.URI,
		Database: cfg.Mongo.Database,
		Timeout:  cfg.Mongo.Timeout,
	})
}

func initMongoFeedRepository(client *mongo.Client) *mongo.FeedRepository {
	return mongo.NewFeedRepository(client)
}
