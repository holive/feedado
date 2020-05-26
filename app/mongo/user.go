package mongo

import "go.mongodb.org/mongo-driver/mongo"

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(conn *Client) *UserRepository {
	return &UserRepository{
		collection: conn.db.Collection("user"),
	}
}
