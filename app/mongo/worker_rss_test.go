package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/holive/feedado/app/rss"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRssWorkerRepository_Find(t *testing.T) {
	client, err := NewConnectionTest()
	defer client.Disconnect(context.Background())
	require.NoError(t, err)

	rssWorkerRepository := NewRssWorkerRepository(NewDBTest(client))

	err = rssWorkerRepository.collection.Drop(context.Background())
	require.NoError(t, err)

	var r rss.RSS
	source1 := "https://google.com"
	source2 := "https://brasil.com"

	{
		feeds := []*rss.RSS{
			{
				Source: source1,
			},
			{
				Source: source2,
			},
		}

		err = rssWorkerRepository.Create(context.Background(), feeds)
		require.NoError(t, err)
	}

	{
		filter := bson.M{"source": bson.M{"$eq": source1}}
		err = rssWorkerRepository.collection.FindOne(context.Background(), filter).Decode(&r)
		require.NoError(t, err)
		assert.Equal(t, r.Source, source1)
	}

	{
		filter := bson.M{"source": bson.M{"$eq": source2}}
		err = rssWorkerRepository.collection.FindOne(context.Background(), filter).Decode(&r)
		require.NoError(t, err)
		assert.Equal(t, r.Source, source2)
	}
}
