package transport

import (
	"app/internal/models"
	"app/pkg/broker"
	"app/pkg/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/wb-go/wbf/retry"
)

const (
	WorkerCountKey = "workerCount"
)

type ConsumerConfig struct {
	WorkerCount int `env:"WORKER_COUNT" env-default:"5"`
	LogLevel    int `env:"LOG_LEVEL" env-default:"1"`
}

type ServiceInterface interface {
	ProcessTask(*models.KafkaTask) error
}

type Consumer struct {
	service ServiceInterface
	ctx     context.Context
	brk     *broker.Broker
	ch      chan kafka.Message
}

func New(service ServiceInterface, cfg *ConsumerConfig, ctx context.Context, brk *broker.Broker) *Consumer {
	return &Consumer{
		service: service,
		// ctx:     context.WithValue(context.WithValue(ctx, WorkerCountKey, cfg.WorkerCount), logger.LoggerKey, logger.LoggerFromCtx(ctx).LoggerLevel(cfg.LogLevel)),
		ctx: context.WithValue(ctx, WorkerCountKey, cfg.WorkerCount),
		brk: brk,
		ch:  make(chan kafka.Message),
	}
}

func (c *Consumer) Consume() {
	lg := logger.LoggerFromCtx(c.ctx).Lg

	for msg := range c.ch {
		var task models.KafkaTask
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			lg.Error().Err(err).Msg("failed to unmarshal Kafka message")
			continue
		}

		if err := c.service.ProcessTask(&task); err != nil {
			lg.Error().Err(err).Str("task_id", task.ID).Msg("failed to process task")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := retry.DoContext(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 1}, func() error {
			ctxPerCommit, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return c.brk.Consumer.Commit(ctxPerCommit, msg)
		})

		if err != nil {
			lg.Error().Err(err).Str("task_id", task.ID).Msg("failed to commit Kafka message")
		} else {
			lg.Debug().Str("task_id", task.ID).Msg("task processed successfully")
		}
	}
}

func (c *Consumer) Start() {

	lg := logger.LoggerFromCtx(c.ctx).Lg

	for range c.ctx.Value(WorkerCountKey).(int) {
		go c.Consume()
	}

	lg.Info().Msg("consumer starting...")
	c.brk.Consumer.StartConsuming(c.ctx, c.ch, retry.Strategy{Attempts: 3, Delay: 2, Backoff: 1.5})
}

func (c *Consumer) Stop() {
	lg := logger.LoggerFromCtx(c.ctx).Lg

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := retry.DoContext(ctx, retry.Strategy{Attempts: 3, Delay: 2, Backoff: 1.5}, c.brk.Consumer.Close)
	if err != nil {
		lg.Error().Err(err).Msg("failed to close consumer")
	} else {
		lg.Info().Msg("consumer closed")
	}
}
