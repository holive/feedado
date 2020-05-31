package feedado

import (
	"github.com/holive/feedado/app/feed"
	"github.com/holive/feedado/app/mongo"
	infraHTTP "github.com/holive/gopkg/net/http"
)

func initFeedService(db *mongo.Client, client infraHTTP.Runner) *feed.Service {
	repository := initMongoFeedRepository(db)

	return feed.NewService(repository, client)
}

func initFeedWorkerService(db *mongo.Client, client infraHTTP.Runner) *feed.WorkerService {
	repository := initMongoFeedRepository(db)

	return feed.NewWorkerService(repository, client)
}
