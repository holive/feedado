package config

import (
	"os"
	"time"

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

type Router struct {
	MiddlewareTimeout time.Duration
}

type HTTPClient struct {
	Timeout   time.Duration
	UserAgent string
}

type Config struct {
	HTTPServer *HTTPServer
	HTTPClient *HTTPClient
}

func loadConfig(profile string) error {
	viper.AddConfigPath(".config")
	viper.SetConfigName(profile)

	if err := viper.MergeInConfig(); err != nil {
		return errors.Wrap(err, "can't read the config file")
	}

	return nil
}

func New() (*Config, error) {
	profile := os.Getenv("APP_PROFILE")

	if profile == "" {
		profile = "development"
	}

	if err := loadConfig(profile); err != nil {
		return nil, errors.Wrap(err, "couldn't initialize app config")
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal config file")
	}

	return &c, nil
}
