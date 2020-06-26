package feedado

import (
	"os"

	"go.uber.org/zap"
)

func initLogger() (*zap.SugaredLogger, error) {
	var (
		logger  *zap.Logger
		err     error
		profile = os.Getenv("APP_PROFILE")
	)

	if profile == "production" {
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	}

	logger, err = zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
