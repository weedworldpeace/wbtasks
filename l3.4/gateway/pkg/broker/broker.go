package broker

import (
	"fmt"

	"github.com/wb-go/wbf/kafka"
)

type BrokerConfig struct {
	BrokerHost  string `env:"KAFKA_HOST" env-default:"localhost"`
	BrokerPort  string `env:"KAFKA_PORT" env-default:"9092"`
	BrokerTopic string `env:"KAFKA_TOPIC" env-default:"tasks"`
}

type Broker struct {
	Producer *kafka.Producer
}

func New(cfg BrokerConfig) *Broker {
	return &Broker{
		Producer: kafka.NewProducer([]string{fmt.Sprintf("%s:%s", cfg.BrokerHost, cfg.BrokerPort)}, cfg.BrokerTopic),
	}
}
