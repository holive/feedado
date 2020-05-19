package feed

import (
	infraHTTP "github.com/holive/gopkg/net/http"
	"go.uber.org/zap"
)

type Service struct {
	repo       Repository
	httpRunner infraHTTP.Runner
	logger     *zap.SugaredLogger
}

func NewService(repository Repository, client infraHTTP.Runner, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:       repository,
		httpRunner: client,
		logger:     logger,
	}
}
