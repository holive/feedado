package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/rss"
	"github.com/holive/feedado/app/worker"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
)

type Worker struct {
	Cfg      *config.Config
	Services *WorkerServices
	Worker   *worker.Worker
	runner   infraHTTP.Runner
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

	w.runner, err = initHTTPClient(w.Cfg)
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

	w.Services = w.initWorkerServices()

	err = w.initRssWorker(logger, db, receiver)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Worker) initWorkerServices() *WorkerServices {
	rssService := initRssWorkerService(w.runner)

	return &WorkerServices{
		RSS: rssService,
	}
}
