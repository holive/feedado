package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/rss"
	infraHTTP "github.com/holive/gopkg/net/http"
	"go.uber.org/zap"
)

func initRssService(db *mongo.Client, runner infraHTTP.Runner) *rss.Service {
	repository := initMongoRssRepository(db)

	return rss.NewService(repository, runner)
}

func initRssProcessor(cfg *config.Config,
	logger *zap.SugaredLogger,
	runner infraHTTP.Runner,
	db *mongo.Client) (*rss.Processor, error) {

	rssRepo := initMongoRssWorkerRepository(db)
	schemaGetter := mongo.NewFeedWorkerRepository(db)

	return rss.NewProcessor(rssRepo, cfg.RSSProcessor, runner, logger, schemaGetter)
}

func initRssWorkerService(runner infraHTTP.Runner) *rss.WorkerService {
	return rss.NewWorkerService(runner)
}
