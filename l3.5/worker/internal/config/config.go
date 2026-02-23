package config

import (
	"fmt"
	"worker/internal/models"
	"worker/internal/service"
	"worker/pkg/data"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/multierr"
)

type Config struct {
	ServiceConfig service.ServiceConfig
	DataConfig    data.DataConfig
}

func (c *Config) valid() error {
	var resultErr error

	resultErr = multierr.Append(resultErr, validInterval(c.ServiceConfig.ServiceInterval))
	resultErr = multierr.Append(resultErr, validPort(c.DataConfig.DbPort))

	return resultErr
}

func validInterval(interval int) error {
	if interval <= 0 {
		return models.ErrInvalidInterval
	}
	return nil
}

func validPort(port int) error {
	if port <= 0 || port > 65535 {
		return models.ErrInvalidPort
	}
	return nil
}

func New() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(fmt.Sprintf("failed to read config from env: %v", err))
	}

	if err := cfg.valid(); err != nil {
		panic(fmt.Sprintf("failed to valid config: %v", err))
	}

	return cfg
}
