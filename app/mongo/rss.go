package mongo

import (
	"context"

	"github.com/holive/feedado/app/rss"
	"go.mongodb.org/mongo-driver/mongo"
)

type RSSRepository struct {
	collection *mongo.Collection
}

func (rr *RSSRepository) Create(ctx context.Context, rss *rss.RSS) (*rss.RSS, error) {
	panic("implement me")
}

func NewRssRepository(conn *Client) *RSSRepository {
	return &RSSRepository{
		collection: conn.db.Collection("rss"),
	}
}
