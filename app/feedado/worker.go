package feedado

import (
	"database/sql"

	infraHTTP "github.com/holive/gopkg/net/http"
	"gitlab.vpc-zoom-01/squad-data/coffe/app/config"
	"gitlab.vpc-zoom-01/squad-data/coffe/app/worker"
	"go.uber.org/zap"
)

func initWorkerOffer(cfg config.Config, logger *zap.SugaredLogger, db *sql.DB, runner infraHTTP.Runner) (*worker.Worker, error) {

	processor, err := initOfferProcessor(cfg, logger, db, runner)
	if err != nil {
		return nil, err
	}

	receiver, err := initGoCloudOfferReceiver(cfg)
	if err != nil {
		return nil, err
	}

	return worker.New(cfg.OfferWorker, logger, receiver, processor)
}
