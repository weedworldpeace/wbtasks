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
	ConsumerConfig transport.ConsumerConfig
	DataConfig     data.DataConfig
	BrokerConfig   broker.BrokerConfig
	StorageConfig  storage.StorageConfig
}

func (c *Config) valid() error {
	var resultErr error
	if err := validPort(c.DataConfig.DbPort); err != nil {
		resultErr = multierr.Append(resultErr, err)
	}
	if err := validWorkerCount(c.ConsumerConfig.WorkerCount); err != nil {
		resultErr = multierr.Append(resultErr, err)
	}
	if err := validLogLevel(c.ConsumerConfig.LogLevel); err != nil {
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

func validWorkerCount(workerCount int) error {
	if workerCount < 0 || workerCount > 50 {
		return models.ErrInvalidWorkerCount
	}
	return nil
}

func validReleaseMode(mode string) error {
	if mode != "" && mode != "release" {
		return models.ErrInvalidReleaseMode
	}
	return nil
}

func validLogLevel(level int) error {
	if level < 1 || level > 8 {
		return models.ErrInvalidLogLevel
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
