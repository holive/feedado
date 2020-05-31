package feedado

import "github.com/holive/feedado/app/config"

func loadConfig(configPath string) (*config.Config, error) {
	return config.New(configPath)
}
