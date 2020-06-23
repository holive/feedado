package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/holive/feedado/app/feed"
	"github.com/stretchr/testify/require"
)

func TestFeedRepository_Create(t *testing.T) {
	client, err := NewConnectionTest()
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	feedRepository := NewFeedRepository(NewDBTest(client))
	err = feedRepository.collection.Drop(context.Background())
	require.NoError(t, err)

	f := feed.Feed{
		Source:      "asdf",
		Description: "1234",
		Sections: []feed.Section{
			{
				SectionSelector:  "sdfg",
				TitleSelector:    "xcvb",
				SubtitleSelector: "erty",
				UrlSelector:      "rtyu",
			},
		},
	}

	_, err = feedRepository.Create(context.Background(), &f)
	require.NoError(t, err)

	var result feed.Feed
	err = feedRepository.collection.FindOne(context.Background(), bson.M{"source": "asdf"}).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, result, f)
}

func TestFeedRepository_Update(t *testing.T) {
	client, err := NewConnectionTest()
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	feedRepository := NewFeedRepository(NewDBTest(client))
	err = feedRepository.collection.Drop(context.Background())
	require.NoError(t, err)

	{
		result, err := feedRepository.collection.InsertOne(
			context.Background(),
			bson.D{
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

	f := feed.Feed{
		Source:      "https://google.com",
		Description: "1234",
		Sections: []feed.Section{
			{
				SectionSelector:  "sdfg",
				TitleSelector:    "xcvb",
				SubtitleSelector: "erty",
				UrlSelector:      "rtyu",
			},
		},
	}

	err = feedRepository.Update(context.Background(), &f)
	require.NoError(t, err)

	var result feed.Feed
	err = feedRepository.collection.FindOne(context.Background(), bson.M{"source": "https://google.com"}).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, result.Description, f.Description)
}

func TestFeedRepository_Delete(t *testing.T) {
	client, err := NewConnectionTest()
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	feedRepository := NewFeedRepository(NewDBTest(client))
	err = feedRepository.collection.Drop(context.Background())
	require.NoError(t, err)

	{
		result, err := feedRepository.collection.InsertOne(
			context.Background(),
			bson.D{
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

	err = feedRepository.Delete(context.Background(), "https://google.com")
	require.NoError(t, err)

	singleResult := feedRepository.collection.FindOne(context.Background(), bson.M{"source": "https://google.com"})
	require.Error(t, singleResult.Err())
}

func TestFeedRepository_FindBySource(t *testing.T) {
	client, err := NewConnectionTest()
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	feedRepository := NewFeedRepository(NewDBTest(client))
	err = feedRepository.collection.Drop(context.Background())
	require.NoError(t, err)

	{
		result, err := feedRepository.collection.InsertOne(
			context.Background(),
			bson.D{
				{"source", "https://google.com"},
				{"description", "1234"},
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

	result, err := feedRepository.FindBySource(context.Background(), "https://google.com")
	require.NoError(t, err)

	assert.Equal(t, result.Description, "1234")
}

func TestFeedRepository_FindAll(t *testing.T) {
	client, err := NewConnectionTest()
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	feedRepository := NewFeedRepository(NewDBTest(client))
	err = feedRepository.collection.Drop(context.Background())
	require.NoError(t, err)

	r, err := feedRepository.collection.InsertOne(
		context.Background(),
		bson.D{
			{"source", "https://google.com"},
			{"description", "1234"},
			{"sections", bson.A{}},
		})
	require.NoError(t, err)
	require.NotNil(t, r.InsertedID)

	r, err = feedRepository.collection.InsertOne(
		context.Background(),
		bson.D{
			{"source", "https://brasil.com.br"},
			{"description", "asdf"},
			{"sections", bson.A{}},
		})
	require.NoError(t, err)
	require.NotNil(t, r.InsertedID)

	result, err := feedRepository.FindAll(context.Background(), "", "")
	require.NoError(t, err)

	assert.Equal(t, 2, len(result.Feeds))
	assert.True(t, result.Feeds[0].Description == "1234" || result.Feeds[0].Description == "asdf")
}
