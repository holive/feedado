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

func initMongoUserRepository(client *mongo.Client) *mongo.UserRepository {
	return mongo.NewUserRepository(client)
}

func initMongoRSSRepository(client *mongo.Client) *mongo.RSSRepository {
	return mongo.NewRSSRepository(client)
}
