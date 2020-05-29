package feedado

import (
	"github.com/holive/feedado/app/config"
	"github.com/holive/feedado/app/mongo"
	"github.com/holive/feedado/app/worker"
	infraHTTP "github.com/holive/gopkg/net/http"
	"go.uber.org/zap"
)

func initWorkerRSS(cfg *config.Config, logger *zap.SugaredLogger, db *mongo.Client, runner infraHTTP.Runner) (*worker.Worker, error) {

	processor, err := initRSSProcessor(cfg, logger, db, runner)
	if err != nil {
		return nil, err
	}

	receiver, err := initGoCloudOfferReceiver(cfg)
	if err != nil {
		return nil, err
	}

	return worker.New(cfg.OfferWorker, logger, receiver, processor)
}
