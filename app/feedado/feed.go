package feedado

import (
	"github.com/holive/feedado/app/feed"
	"github.com/holive/feedado/app/mongo"
	infraHTTP "github.com/holive/gopkg/net/http"
	"go.uber.org/zap"
)

func initFeedService(db *mongo.Client, logger *zap.SugaredLogger, client infraHTTP.Runner) *feed.Service {
	repository := initMongoFeedRepository(db)

	return feed.NewService(repository, client, logger)
}
