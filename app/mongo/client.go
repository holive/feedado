package mongo

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	db *mongo.Client
}

type ClientConfig struct {
	URI      string
	User     string
	Password string
	AppName  string
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
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to database")
	}

	return &Client{
		db: client,
	}, nil
}
