package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/feed"
	"github.com/holive/feedado/app/mongo"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Feedado struct {
	Cfg      *config.Config
	Services *Services
}

type Services struct {
	Feed *feed.Service
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

	f.Services, err = initServices(f.Cfg, db, httpClient, logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize services")
	}

	return f, nil
}

func initServices(cfg *config.Config, db *mongo.Client, client infraHTTP.Runner, logger *zap.SugaredLogger) (*Services, error) {
	feedService := initFeedService(db, logger, client)

	return &Services{
		Feed: feedService,
	}, nil
}
