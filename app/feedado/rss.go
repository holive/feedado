package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/gocloud"
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/rss"
	"github.com/holive/feedado/app/worker"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func initRssService(db *mongo.Client, runner infraHTTP.Runner) *rss.Service {
	repository := initMongoRssRepository(db)

	return rss.NewService(repository, runner)
}

func initRssWorkerService(runner infraHTTP.Runner) *rss.WorkerService {
	return rss.NewWorkerService(runner)
}

func initRssProcessor(cfg *config.Config, logger *zap.SugaredLogger, runner infraHTTP.Runner,
	db *mongo.Client) (*rss.Processor, error) {

	updater := initMongoRssWorkerRepository(db)
	schemaGetter := mongo.NewFeedWorkerRepository(db)

	return rss.NewProcessor(updater, cfg.RSSProcessor, runner, logger, schemaGetter)
}

func (w *Worker) initRssWorker(logger *zap.SugaredLogger, db *mongo.Client, receiver *gocloud.RSSReceiver) error {
	processor, err := initRssProcessor(w.Cfg, logger, w.runner, db)
	if err != nil {
		return errors.Wrap(err, "could not initialize worker rss processor")
	}

	w.Worker, err = worker.New(w.Cfg.RSSWorker, logger, receiver, processor)
	if err != nil {
		return errors.Wrap(err, "could not initialize the worker")
	}

	return nil
}
