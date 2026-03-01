package config

import (
	"app/internal/models"
	"app/internal/service"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/multierr"
)

type Config struct {
	ServiceConfig service.ServiceConfig
}

func (c *Config) valid() error {
	var resultErr error

	resultErr = multierr.Append(resultErr, validWrkCount(c.ServiceConfig.WrkCount))

	for i := range c.ServiceConfig.WrkCount {
		resultErr = multierr.Append(resultErr, validPort(c.ServiceConfig.Wrks[i].Port))
	}

	return resultErr
}

func validPort(port int) error {
	if port <= 0 || port > 65535 {
		return models.ErrInvalidPort
	}
	return nil
}

func validWrkCount(count int) error {
	if count < 1 || count > 20 {
		return models.ErrInvalidWrkCount
	}
	return nil
}

func New() *Config {
	cfg := Config{}

	if err := cleanenv.ReadConfig("./config.yaml", &cfg.ServiceConfig); err != nil {
		panic(fmt.Sprintf("failed to read config from file: %v", err))
	}

	if err := cfg.valid(); err != nil {
		panic(fmt.Sprintf("failed to valid config: %v", err))
	}

	return &cfg
}
