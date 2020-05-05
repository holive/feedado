package feed

import "github.com/holive/feed/app/config"

func loadConfig() (*config.Config, error) {
	return config.New()
}
