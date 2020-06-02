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

func initRSSProcessor(cfg *config.Config, logger *zap.SugaredLogger, repository rss.Updater, runner infraHTTP.Runner) (*rss.Processor, error) {
	return rss.NewProcessor(repository, cfg.RSSProcessor, runner, logger)
}

func initRssWorkerService(db *mongo.Client, runner infraHTTP.Runner) *rss.WorkerService {
	repository := initMongoRssWorkerRepository(db)

	return rss.NewWorkerService(repository, runner)
}
