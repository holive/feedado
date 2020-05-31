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

type FeedadoWorker struct {
	Cfg      *config.Config
	Services *WorkerServices
}

type WorkerServices struct {
	Feed   *feed.Service
	RSS    *rss.Service
	worker *worker.Worker
}

func NewWorker() (*FeedadoWorker, error) {
	var (
		err error
		w   = &FeedadoWorker{}
	)

	w.Cfg, err = loadConfig("./config/worker")
	if err != nil {
		return nil, errors.Wrap(err, "could not load config")
	}

	db, err := initMongoClient(w.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize mongo client")
	}

	httpClient, err := initHTTPClient(w.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize http client")
	}

	logger, err := initLogger()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize logger")
	}

	w.Services, err = initWorkerServices(w.Cfg, db, httpClient, logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize worker services")
	}

	return w, nil
}

func initWorkerServices(cfg *config.Config, db *mongo.Client, client infraHTTP.Runner, logger *zap.SugaredLogger) (
	*WorkerServices, error) {
	feedService := initFeedService(db, client)
	rssService := initRssService(db, client)

	processor, err := initRSSProcessor(cfg, logger, rssService, client)
	if err != nil {
		return nil, err
	}

	receiver, err := initGoCloudOfferReceiver(cfg)
	if err != nil {
		return nil, err
	}

	w, err := worker.New(cfg.RSSWorker, logger, receiver, processor)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize the worker")
	}

	return &WorkerServices{
		Feed:   feedService,
		RSS:    rssService,
		worker: w,
	}, nil
}
