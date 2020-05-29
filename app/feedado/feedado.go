package feedado

import (
	"github.com/holive/feedado/app/worker"

	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/feed"
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/user"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Feedado struct {
	Cfg        *config.Config
	Services   *Services
	FeedWorker *worker.Worker
}

type Services struct {
	Feed *feed.Service
	User *user.Service
}

func New() (*Feedado, error) {
	var (
		err error
		f   = &Feedado{}
	)

	f.Cfg, err = loadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "could not load config")
	}

	db, err := initMongoClient(f.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize mongo client")
	}

	httpClient, err := initHTTPClient(f.Cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize http client")
	}

	logger, err := initLogger()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize logger")
	}

	f.Services, err = initServices(db, httpClient, logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize services")
	}

	f.FeedWorker, err = initWorkerRSS(f.Cfg, logger, db, httpClient)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize the worker")
	}

	return f, nil
}

func initServices(db *mongo.Client, client infraHTTP.Runner, logger *zap.SugaredLogger) (*Services, error) {
	feedService := initFeedService(db, client)
	userService := initUserService(db)

	return &Services{
		Feed: feedService,
		User: userService,
	}, nil
}
