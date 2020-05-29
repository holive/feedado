package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/gocloud"
)

func initGoCloudRSSPublisher(cfg *config.Config) (*gocloud.RSSPublisher, error) {
	return gocloud.NewOfferPublisher(cfg.RSSPubSub)
}

func initGoCloudOfferReceiver(cfg *config.Config) (*gocloud.RSSReceiver, error) {
	return gocloud.NewOfferReceiver(cfg.RSSPubSub)
}
