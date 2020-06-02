package rss

import (
	"context"
	"encoding/json"

	"github.com/holive/feedado/app/feed"

	"go.uber.org/zap"

	infraHTTP "github.com/holive/gopkg/net/http"

	"github.com/pkg/errors"
)

type Processor struct {
	updater      WorkerRepository
	schemaGetter SchemaGetter
	userAgent    string
	runner       infraHTTP.Runner
	logger       *zap.SugaredLogger
}

type ProcessorConfig struct {
	UserAgent string
}

func (w *Processor) Process(ctx context.Context, message []byte) error {
	var m map[string]string
	if err := json.Unmarshal(message, &m); err != nil {
		return errors.Wrap(err, "could not unmarshal message")
	}

	//id, err := primitive.ObjectIDFromHex(m["_id"])
	//if err != nil {
	//	return err
	//}

	var schema *feed.Feed
	schema, err := w.schemaGetter.Find(ctx, m["_id"])
	if err != nil {
		return errors.Wrap(err, "could not find schema")
	}

	// TODO: fetch webpage info
	_ = schema

	return nil
}

func NewProcessor(updater WorkerRepository, cfg *ProcessorConfig, runner infraHTTP.Runner, logger *zap.SugaredLogger) (*Processor, error) {
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
