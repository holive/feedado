package feedado

import "github.com/holive/feedado/app/config"

func loadConfig() (*config.Config, error) {
	return config.New()
}
