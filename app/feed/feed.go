package feed

import (
	"github.com/holive/feed/app/config"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Feed struct {
	Cfg      *config.Config
	Services *Services
}

type Services interface {
}

func New() (*Feed, error) {
	var (
		err error
		f   = &Feed{}
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

func initServices(cfg *config.Config, client infraHTTP.Runner, logger *zap.SugaredLogger) (*Services, error) {

	return nil, nil
}
