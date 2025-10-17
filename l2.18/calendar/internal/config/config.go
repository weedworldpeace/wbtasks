package config

import (
	"calendar/internal/transport"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	transport.ServerConfig
}

func New() *Config {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to extract logger: %v", err))
	}

	return cfg
}
