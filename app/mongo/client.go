package mongo

import (
	"context"
	"time"

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
	// write timeout ?
	// poolSize ?
	// autoReconnect ?
	// noDelay ?
	// keepAlive ?
	// connectTimeoutMS ?
	// loggerLevel ?
	// logger ?
}

func New(cfg *ClientConfig) (*Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), cfg.Timeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to database")
	}

	return &Client{
		db: client,
	}, nil
}
