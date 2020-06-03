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
	updater      Updater
	schemaGetter SchemaGetter
	userAgent    string
	runner       infraHTTP.Runner
	logger       *zap.SugaredLogger
}

type ProcessorConfig struct {
	UserAgent string
}

func (w *Processor) Process(ctx context.Context, message []byte) error {
	var m struct {
		ID string `json:"_id"`
	}

	if err := json.Unmarshal(message, &m); err != nil {
		return errors.Wrap(err, "could not unmarshal message")
	}

	//id, err := primitive.ObjectIDFromHex(m["_id"])
	//if err != nil {
	//	return err
	//}

	schema, err := w.schemaGetter.Find(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "could not find schema")
	}

	rssResults, err := fetchRssResults(schema)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal message")
	}

	err = w.updater.Create(ctx, rssResults)

	return nil
}

func fetchRssResults(schema *feed.Feed) ([]*RSS, error) {
	panic("implement me")
}

func NewProcessor(updater Updater,
	cfg *ProcessorConfig,
	runner infraHTTP.Runner,
	logger *zap.SugaredLogger,
	schemaGetter SchemaGetter) (*Processor, error) {

	if updater == nil {
		return nil, errors.New("updater can't be nil")
	}

	if schemaGetter == nil {
		return nil, errors.New("schemaGetter can't be nil")
	}

	if cfg.UserAgent == "" {
		return nil, errors.New("config can't has empty fields")
	}

	return &Processor{
		updater:      updater,
		schemaGetter: schemaGetter,
		userAgent:    cfg.UserAgent,
		runner:       runner,
		logger:       logger,
	}, nil
}
