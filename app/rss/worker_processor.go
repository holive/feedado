package rss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

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

func (p *Processor) Process(ctx context.Context, message []byte) error {
	fmt.Println("Process number :" + string(message))
	var m struct {
		SchemaID string `json:"schema_id"`
	}

	if err := json.Unmarshal(message, &m); err != nil {
		return errors.Wrap(err, "could not unmarshal message")
	}

	schema, err := p.schemaGetter.Find(ctx, m.SchemaID)
	if err != nil {
		return errors.Wrap(err, "could not find schema")
	}

	rssResults, err := p.fetchRssResults(schema)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal message")
	}

	err = p.updater.Create(ctx, rssResults)

	return nil
}

func (p *Processor) fetchRssResults(schema *feed.Feed) ([]*RSS, error) {
	req, err := http.NewRequest("GET", schema.Source, nil)
	if err != nil {
		return nil, err
	}

	res, err := p.runner.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch "+schema.Source)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.Wrapf(err, "status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	selector := strings.Trim(schema.Sections[0].ParentBlockClass+" "+schema.Sections[0].EachBlockClass, " ")

	var title, subtitle string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		title = s.Find(schema.Sections[0].Title).Text()
		subtitle = s.Find(schema.Sections[0].Subtitle).Text()
		fmt.Printf("title subtitle %d: %s - %s\n", i, title, subtitle)
	})

	//body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))

	return nil, nil
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
