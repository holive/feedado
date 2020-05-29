package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type RSSRepository struct {
	collection *mongo.Collection
}

func NewRSSRepository(conn *Client) *RSSRepository {
	return &RSSRepository{
		collection: conn.db.Collection("rss"),
	}
}
