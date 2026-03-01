package config

import (
	"app/internal/models"
	"app/internal/transport"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/multierr"
)

type Config struct {
	ServerConfig transport.ServerConfig
}

func (c *Config) valid() error {
	var resultErr error

	resultErr = multierr.Append(resultErr, validPort(c.ServerConfig.Port))
	resultErr = multierr.Append(resultErr, validReleaseMode(c.ServerConfig.ReleaseMode))

	return resultErr
}

func validPort(port int) error {
	if port < 0 || port > 65535 {
		return models.ErrInvalidPort
	}
	return nil
}

func validReleaseMode(mode string) error {
	if mode != "" && mode != "release" {
		return models.ErrInvalidReleaseMode
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
