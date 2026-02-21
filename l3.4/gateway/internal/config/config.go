package config

import (
	"app/internal/models"
	"app/internal/transport"
	"app/pkg/broker"
	"app/pkg/data"
	"app/pkg/storage"
	"fmt"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/multierr"
)

type Config struct {
	ServerConfig  transport.ServerConfig
	DataConfig    data.DataConfig
	BrokerConfig  broker.BrokerConfig
	StorageConfig storage.StorageConfig
}

func (c *Config) valid() error {
	var resultErr error
	if err := validPort(c.ServerConfig.Port); err != nil {
		resultErr = multierr.Append(resultErr, err)
	}
	if err := validPort(c.DataConfig.DbPort); err != nil {
		resultErr = multierr.Append(resultErr, err)
	}
	if err := validReleaseMode(c.ServerConfig.ReleaseMode); err != nil {
		resultErr = multierr.Append(resultErr, err)
	}
	return resultErr
}

func validPort(port string) error {
	convPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if convPort < 0 || convPort > 65535 {
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
