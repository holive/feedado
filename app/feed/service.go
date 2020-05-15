package feed

import (
	"github.com/holive/feedado/app/mongo"
	infraHTTP "github.com/holive/gopkg/net/http"
	"go.uber.org/zap"
)

type Service struct {
	repo       Repository
	httpRunner infraHTTP.Runner
	logger     *zap.SugaredLogger
}

func NewService(repository *mongo.FeedRepository, client infraHTTP.Runner, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
		logger:     logger,
	}
}
