package rss

import (
	"context"

	"go.uber.org/zap"

	infraHTTP "github.com/holive/gopkg/net/http"

	"github.com/pkg/errors"
)

type Processor struct {
	updater   Repository
	userAgent string
	runner    infraHTTP.Runner
	logger    *zap.SugaredLogger
}

type ProcessorConfig struct {
	UserAgent string
}

func (w *Processor) Process(ctx context.Context, message []byte) error {
	panic("missing Process implementation")
}

func NewProcessor(updater Repository, cfg *ProcessorConfig, runner infraHTTP.Runner, logger *zap.SugaredLogger) (*Processor, error) {
	if updater == nil {
		return nil, errors.New("updater can't be nil")
	}

	if cfg.UserAgent == "" {
		return nil, errors.New("config can't has empty fields")
	}

	return &Processor{
		updater:   updater,
		userAgent: cfg.UserAgent,
		runner:    runner,
		logger:    logger,
	}, nil
}
