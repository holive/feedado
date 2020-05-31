package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/feed"
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/rss"
	"github.com/holive/feedado/app/worker"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Worker struct {
	Cfg      *config.Config
	Services *WorkerServices
	worker   *worker.Worker
	runner   infraHTTP.Runner
}

type WorkerServices struct {
	Feed *feed.WorkerService
	RSS  *rss.WorkerService
}

func NewWorker() (*Worker, error) {
	var (
		err error
		w   = &Worker{}
	)

	w.Cfg, err = loadConfig("./config/worker")
	if err != nil {
		return nil, errors.Wrap(err, "could not load config")
	}

	db, err := initMongoClient(w.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize mongo client")
	}

	w.runner, err = initHTTPClient(w.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize http client")
	}

	logger, err := initLogger()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize logger")
	}

	w.Services = w.initWorkerServices(db)

	err = w.initWorker(logger)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Worker) initWorkerServices(db *mongo.Client) *WorkerServices {
	feedService := initFeedWorkerService(db, w.runner)
	rssService := initRssWorkerService(db, w.runner)

	return &WorkerServices{
		Feed: feedService,
		RSS:  rssService,
	}
}

func (w *Worker) initWorker(logger *zap.SugaredLogger) error {
	processor, err := initRSSProcessor(w.Cfg, logger, w.Services.RSS, w.runner)
	if err != nil {
		return errors.Wrap(err, "could not initialize worker rss processor")
	}

	receiver, err := initGoCloudOfferReceiver(w.Cfg)
	if err != nil {
		return errors.Wrap(err, "could not initialize worker rss receiver")
	}

	w.worker, err = worker.New(w.Cfg.RSSWorker, logger, receiver, processor)
	if err != nil {
		return errors.Wrap(err, "could not initialize the worker")
	}

	return nil
}
