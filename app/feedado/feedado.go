package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/feed"
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/rss"
	"github.com/holive/feedado/app/user"
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
	User *user.Service
	RSS  *rss.Service
}

func New() (*Feedado, error) {
	var (
		err error
		f   = &Feedado{}
	)

	f.Cfg, err = loadConfig("./config/api")
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

	f.Services = initServices(f.Cfg, db, httpClient, logger)

	return f, nil
}

func initServices(cfg *config.Config, db *mongo.Client, client infraHTTP.Runner, logger *zap.SugaredLogger) *Services {
	feedService := initFeedService(db, client)
	userService := initUserService(db)
	rssService := initRssService(db, client)

	return &Services{
		Feed: feedService,
		User: userService,
		RSS:  rssService,
	}
}
