package gocloud

import (
	"context"
	"encoding/json"

	"github.com/holive/feedado/app/feed"

	"github.com/pkg/errors"

	"gocloud.dev/pubsub"
)

type RSSPubSubCfg struct {
	Name    string
	Service string
	Region  string
}

type RSSPublisher struct {
	client *Client
}

func (op *RSSPublisher) Publish(ctx context.Context, feedId feed.SQS) error {
	rawMessage, err := json.Marshal(feedId)
	if err != nil {
		return errors.Wrap(err, "Cloud not encode update message")
	}

	return op.client.Publish(ctx, rawMessage)
}

func NewRssPublisher(cfg *RSSPubSubCfg) (*RSSPublisher, error) {
	client, err := NewClient(
		SetResourceName(cfg.Name),
		SetService(cfg.Service),
		SetRegion(cfg.Region))

	if err != nil {
		return nil, err
	}

	return &RSSPublisher{client: client}, nil
}

type RSSReceiver struct {
	client *Client
}

func (or *RSSReceiver) Receive(ctx context.Context) (*pubsub.Message, error) {
	return or.client.Receive(ctx)
}

func NewRssReceiver(cfg *RSSPubSubCfg) (*RSSReceiver, error) {
	client, err := NewClient(
		SetResourceName(cfg.Name),
		SetService(cfg.Service),
		SetRegion(cfg.Region))

	if err != nil {
		return nil, err
	}

	return &RSSReceiver{client: client}, nil
}
