package mongo

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConnectionTest() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewDBTest(client *mongo.Client) *Client {
	return &Client{db: client.Database("test")}
}

func TestFeedWorkerRepository_Find(t *testing.T) {
	client, err := NewConnectionTest()
	defer client.Disconnect(context.Background())
	require.NoError(t, err)

	feedWorkerRepository := NewFeedWorkerRepository(NewDBTest(client))

	err = feedWorkerRepository.collection.Drop(context.Background())
	require.NoError(t, err)

	id := "5ed848cbc62f0bbb7b8e411e"

	{
		ObjectId, err := primitive.ObjectIDFromHex(id)
		require.NoError(t, err)
		result, err := feedWorkerRepository.collection.InsertOne(
			context.Background(),
			bson.D{
				{"_id", ObjectId},
				{"source", "https://google.com"},
				{"description", ""},
				{"sections", bson.A{
					bson.D{
						{"parent_block_class", ""},
						{"each_block_class", ""},
						{"title", ""},
						{"subtitle", ""},
						{"url", ""},
					},
				}},
			})

		require.NoError(t, err)
		require.NotNil(t, result.InsertedID)
	}

	{
		_, err = feedWorkerRepository.Find(context.Background(), id)
		require.NoError(t, err)
	}

}
