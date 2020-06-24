package feedado

import (
	"github.com/holive/feedado/app/feed"
	"github.com/holive/feedado/app/mongo"
)

func initFeedService(db *mongo.Client) *feed.Service {
	repository := initMongoFeedRepository(db)

	return feed.NewService(repository)
}
