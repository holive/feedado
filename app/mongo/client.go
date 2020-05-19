package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	db *mongo.Client
}

type ClientConfig struct {
	URI     string
	AppName string
	Timeout time.Duration
}

func New(cfg *ClientConfig) (*Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), cfg.Timeout*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to database")
	}

	return &Client{
		db: client,
	}, nil
}