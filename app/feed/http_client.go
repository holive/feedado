package feed

import (
	"net/http"

	"github.com/holive/feed/app/config"
)

func initHTTPClient(config *config.Config) (http.Client, error) {

	return http.Client{}, nil
}
