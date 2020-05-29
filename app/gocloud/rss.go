package gocloud

import (
	"context"
	"encoding/json"

	"github.com/holive/feedado/app/rss"
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

func (op *RSSPublisher) Publish(ctx context.Context, r rss.RSS) error {

	rawMessage, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "Cloud not encode update message")
	}

	return op.client.Publish(ctx, rawMessage)
}

func NewOfferPublisher(cfg *RSSPubSubCfg) (*RSSPublisher, error) {
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

func NewOfferReceiver(cfg *RSSPubSubCfg) (*RSSReceiver, error) {
	client, err := NewClient(
		SetResourceName(cfg.Name),
		SetService(cfg.Service),
		SetRegion(cfg.Region))

	if err != nil {
		return nil, err
	}

	return &RSSReceiver{client: client}, nil
}
