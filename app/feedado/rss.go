package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/rss"
	infraHTTP "github.com/holive/gopkg/net/http"
	"go.uber.org/zap"
)

func initRssService(db *mongo.Client, runner infraHTTP.Runner) *rss.Service {
	repository := initMongoRSSRepository(db)
	rssService := rss.NewService(repository, runner)

	return rssService
}

func initRSSProcessor(cfg *config.Config, logger *zap.SugaredLogger, repository rss.Repository, runner infraHTTP.Runner) (*rss.Processor, error) {
	return rss.NewProcessor(repository, cfg.RSSProcessor, runner, logger)
}
