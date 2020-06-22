package rss

import (
	infraHTTP "github.com/holive/gopkg/net/http"
)

type WorkerService struct {
	httpRunner infraHTTP.Runner
}

//func (w WorkerService) Create(ctx context.Context, feeds []*RSS) error {
//	// não precisa do create aqui.. só um trigger pro
//	panic("implement me")
//}

func NewWorkerService(client infraHTTP.Runner) *WorkerService {
	return &WorkerService{
		httpRunner: client,
	}
}
