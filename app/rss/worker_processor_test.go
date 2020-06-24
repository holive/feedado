package rss

import (
	"context"

	"github.com/holive/feedado/app/feed"
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

//func TestProcessor_fetchRssResults(t *testing.T) {
//	logger := initLoggerTest()
//	updater := UpdaterTest{}
//	schemaGetter := SchemaGetterTest{}
//	cfg := ProcessorConfig{UserAgent: "feedado-test"}
//	runner := &http.Client{}
//
//	processor, err := NewProcessor(&updater, &cfg, runner, logger, &schemaGetter)
//	require.NoError(t, err)
//
//	// request and parse goquery
//	res, err := http.Get("https://economia.estadao.com.br")
//	require.NoError(t, err)
//	defer res.Body.Close()
//	bodyByte, err := ioutil.ReadAll(res.Body)
//	require.NoError(t, err)
//
//	// and request and parse
//
//	rsss, err := processor.sourceResponseToRSS(bodyByte, &feed.Feed{
//		Source:      "https://economia.estadao.com.br",
//		Description: "",
//		Sections: []feed.Section{
//			{
//				SectionSelector:  ".row.management section",
//				TitleSelector:    ".text-wrapper h3",
//				SubtitleSelector: ".text-wrapper p",
//				UrlSelector:      ".text-wrapper a",
//			},
//		},
//	})
//	require.NoError(t, err)
//	assert.Equal(t, 4, len(rsss))
//}
