package rss

import (
	"context"
	netHttp "net/http"
	"testing"

	"github.com/holive/feedado/app/feed"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func initLoggerTest() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

type UpdaterTest struct {
}

type SchemaGetterTest struct {
}

func (u *UpdaterTest) Create(ctx context.Context, feeds []*RSS) error {
	return nil
}

func (s *SchemaGetterTest) Find(ctx context.Context, id string) (*feed.Feed, error) {
	return nil, nil
}

func TestProcessor_fetchRssResults(t *testing.T) {
	logger := initLoggerTest()
	updater := UpdaterTest{}
	schemaGetter := SchemaGetterTest{}
	cfg := ProcessorConfig{UserAgent: "feedado-test"}
	runner := &netHttp.Client{}

	processor, err := NewProcessor(&updater, &cfg, runner, logger, &schemaGetter)
	require.NoError(t, err)

	_, err = processor.fetchRssResults(&feed.Feed{
		Source:      "https://economia.estadao.com.br/",
		Description: "",
		Sections: []feed.Section{
			{
				ParentBlockClass: ".row.management",
				EachBlockClass:   "section",
				Title:            ".text-wrapper h3",
				Subtitle:         ".text-wrapper p",
				Url:              ".text-wrapper",
			},
		},
	})
	require.NoError(t, err)
}

//func TestRssWorkerRepository_Find(t *testing.T) {
//		err = rssWorkerRepository.Create(context.Background(), feeds)
//		require.NoError(t, err)
//	}
//
//	{
//		filter := bson.M{"source": bson.M{"$eq": source1}}
//		err = rssWorkerRepository.collection.FindOne(context.Background(), filter).Decode(&r)
//		require.NoError(t, err)
//		assert.Equal(t, r.Source, source1)
//	}
//
//	{
//		filter := bson.M{"source": bson.M{"$eq": source2}}
//		err = rssWorkerRepository.collection.FindOne(context.Background(), filter).Decode(&r)
//		require.NoError(t, err)
//		assert.Equal(t, r.Source, source2)
//	}
//}
