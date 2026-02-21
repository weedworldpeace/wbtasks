package broker

import (
	"fmt"

	"github.com/wb-go/wbf/kafka"
)

type BrokerConfig struct {
	BrokerHost  string `env:"KAFKA_HOST" env-default:"localhost"`
	BrokerPort  string `env:"KAFKA_PORT" env-default:"9092"`
	BrokerTopic string `env:"KAFKA_TOPIC" env-default:"tasks"`
	BrokerGroup string `env:"KAFKA_GROUP" env-default:"processor-group"`
}

type Broker struct {
	Consumer *kafka.Consumer
}

func New(cfg BrokerConfig) *Broker {
	return &Broker{
		Consumer: kafka.NewConsumer([]string{fmt.Sprintf("%s:%s", cfg.BrokerHost, cfg.BrokerPort)}, cfg.BrokerTopic, cfg.BrokerGroup),
	}
}
