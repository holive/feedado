package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/rss"
	"github.com/holive/feedado/app/worker"
	"github.com/pkg/errors"
)

type Worker struct {
	Cfg      *config.Config
	Services *WorkerServices
	Worker   *worker.Worker
}

type WorkerServices struct {
	RSS *rss.WorkerService
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

	httpRunner, err := initHTTPClient(w.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize http client")
	}

	logger, err := initLogger()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize logger")
	}

	receiver, err := initGoCloudOfferReceiver(w.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize worker rss receiver")
	}

	publisher, err := initGoCloudRSSPublisher(w.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize worker rss publisher")
	}

	w.Services = initRssWorkerService(db, logger, publisher)

	err = w.initRssWorker(logger, db, receiver, httpRunner)
	if err != nil {
		return nil, err
	}

	return w, nil
}
