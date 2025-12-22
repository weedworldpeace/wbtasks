package rabbit

import (
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
)

type RabbitConfig struct {
	Host     string `env:"RABBIT_HOST" env-default:"localhost"`
	Port     string `env:"RABBIT_PORT" env-default:"5672"`
	Username string `env:"RABBIT_USERNAME" env-default:"guest"`
	Password string `env:"RABBIT_PASSWORD" env-default:"guest"`
	Retry    int    `env:"RABBIT_RETRY" env-default:"5"`
	Pause    int    `env:"RABBIT_PAUSE" env-default:"1"`
}

func (rc *RabbitConfig) URL() string {
	return "amqp://" + rc.Username + ":" + rc.Password + "@" + rc.Host + ":" + rc.Port + "/"
}

func New(cfg RabbitConfig) (*rabbitmq.Connection, *rabbitmq.Channel) {
	fmt.Println(cfg.URL())
	conn, err := rabbitmq.Connect(cfg.URL(), cfg.Retry, time.Duration(cfg.Pause)*time.Second)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	if err = ch.ExchangeDeclare("main_exchange", "x-delayed-message", true, false, false, false, amqp091.Table{"x-delayed-type": "direct"}); err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare("main_queue", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	if err = ch.QueueBind(q.Name, "main_routing_key", "main_exchange", false, nil); err != nil {
		panic(err)
	}

	return conn, ch
}
