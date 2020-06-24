package rss

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/holive/feedado/app/feed"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

func (p *Processor) Process(ctx context.Context, message []byte) error {
	var m feed.FeedSQS

	if err := json.Unmarshal(message, &m); err != nil {
		return errors.Wrap(err, "could not unmarshal message")
	}

	schema, err := p.schemaGetter.Find(ctx, m.ID)
	if err != nil {
		return errors.Wrap(err, "could not find schema")
	}

	doc, err := p.fetchSourcePage(schema)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal message")
	}

	rssResults, err := p.sourceResponseToRSS(doc, schema)
	if err != nil {
		return errors.Wrap(err, "could not parse html page to rss")
	}

	err = p.updater.Create(ctx, rssResults)
	if err != nil {
		return errors.Wrapf(err, "could not insert batch of %s", schema.Source)
	}

	p.logger.Info("finishing batch of source: ", schema.Source)

	return nil
}

func (p *Processor) fetchSourcePage(schema *feed.Feed) ([]byte, error) {
	req, err := http.NewRequest("GET", schema.Source, nil)
	if err != nil {
		return nil, err
	}

	response, err := p.runner.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch "+schema.Source)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.Wrapf(err, "status code error: %d %s", response.StatusCode, response.Status)
	}

	return ioutil.ReadAll(response.Body)
}

func (p *Processor) sourceResponseToRSS(htmlPage []byte, schema *feed.Feed) ([]*RSS, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlPage))
	if err != nil {
		p.logger.Errorf("could not find %s", schema.Source)
	}

	rsss := make(map[string]RSS)

	for _, section := range schema.Sections {
		doc.Find(section.SectionSelector).Each(func(i int, s *goquery.Selection) {
			rss := RSS{
				Source:    schema.Source,
				Title:     s.Find(section.TitleSelector).Text(),
				Subtitle:  s.Find(section.SubtitleSelector).Text(),
				URL:       s.Find(section.UrlSelector).AttrOr("href", ""),
				Timestamp: time.Now().Unix(),
			}

			if rss.Title == "" || rss.Source == "" {
				return
			}

			rsss[rss.Title] = rss
		})
	}

	var result []*RSS
	for _, rss := range rsss {
		rss := rss
		result = append(result, &rss)
	}

	return result, nil
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
