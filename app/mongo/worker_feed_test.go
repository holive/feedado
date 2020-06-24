package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/require"
)

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
						{"section_selector", ""},
						{"title_selector", ""},
						{"subtitle_selector", ""},
						{"url_selector", ""},
					},
				}},
			})

		require.NoError(t, err)
		require.NotNil(t, result.InsertedID)
	}

	_, err = feedWorkerRepository.Find(context.Background(), id)
	require.NoError(t, err)
}
