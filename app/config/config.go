package config

import (
	"os"
	"time"

	"github.com/holive/feedado/app/worker"

	"github.com/holive/feedado/app/gocloud"
	"github.com/holive/feedado/app/rss"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	Addr              string
	MaxHeaderBytes    int
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	Router            *Router
}

type DB struct {
	URI      string
	Database string
	Timeout  time.Duration
	RSSttl   int32
}

type Router struct {
	MiddlewareTimeout time.Duration
}

type HTTPClient struct {
	Timeout   time.Duration
	UserAgent string
}

type Config struct {
	HTTPServer   *HTTPServer
	HTTPClient   *HTTPClient
	Mongo        *DB
	RSSPubSub    *gocloud.RSSPubSubCfg
	RSSWorker    *worker.Options
	RSSProcessor *rss.ProcessorConfig
}

func New(configPath string) (*Config, error) {
	profile := os.Getenv("APP_PROFILE")

	if profile == "" {
		profile = "development"
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(profile)

	if err := viper.MergeInConfig(); err != nil {
		return nil, errors.Wrap(err, "can't read the config file")
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal config file")
	}

	return &c, nil
}
