package feed

import (
	"github.com/holive/feed/app/config"
	"github.com/pkg/errors"
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
		t   = &Feed{}
	)

	t.Cfg, err = loadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "could not load config")
	}

	return t, nil
}
