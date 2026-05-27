package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"task-1/internal/config"
	"task-1/internal/model"
)

type RabbitMq struct {
	client    *rabbitmq.RabbitClient
	publisher *rabbitmq.Publisher
}

func New() *RabbitMq {
	url := fmt.Sprintf(
		"amqp://guest:guest@%s%s/",
		config.Cfg.RabbitMq.Host,
		config.Cfg.RabbitMq.Port,
	)

	clientConfig := rabbitmq.ClientConfig{
		URL:            url,
		ConnectionName: "task-1",
		ConnectTimeout: 10 * time.Second,
		Heartbeat:      10 * time.Second,
		ReconnectStrat: retry.Strategy{
			Attempts: 0,
			Delay:    time.Second,
			Backoff:  2,
		},
		ProducingStrat: retry.Strategy{
			Attempts: 3,
			Delay:    time.Second,
			Backoff:  2,
		},
		ConsumingStrat: retry.Strategy{
			Attempts: 3,
			Delay:    time.Second,
			Backoff:  2,
		},
	}

	client, err := rabbitmq.NewClient(clientConfig)
	if err != nil {
		log.Fatal("could not create rabbitmq client: ", err)
	}

	err = client.DeclareExchange("notifications", "direct", true, false, false, nil)
	if err != nil {
		log.Fatal("could not declare exchange: ", err)
	}

	err = client.DeclareQueue(
		"notification",
		"notifications",
		"notification",
		true,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Fatal("could not declare queue: ", err)
	}

	publisher := rabbitmq.NewPublisher(client, "notifications", "application/json")

	return &RabbitMq{
		client:    client,
		publisher: publisher,
	}
}

func (r *RabbitMq) Publish(notification model.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		zlog.Logger.Error().Msg("could not marshal notification to send to rabbitmq: " + err.Error())
		return err
	}

	ctx := context.Background()
	return r.publisher.Publish(ctx, body, "notification")
}

func (r *RabbitMq) Consume(ctx context.Context) (<-chan []byte, error) {
	messages := make(chan []byte, 100)

	messageHandler := func(ctx context.Context, delivery amqp.Delivery) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case messages <- delivery.Body:
			return nil
		}
	}

	consumerConfig := rabbitmq.ConsumerConfig{
		Queue:       "notification",
		ConsumerTag: "task-1-consumer",
		AutoAck:     true,
		Ask: rabbitmq.AskConfig{
			Multiple: false,
		},
		Nack: rabbitmq.NackConfig{
			Multiple: false,
			Requeue:  true,
		},
		Workers:       1,
		PrefetchCount: 1,
	}

	consumer := rabbitmq.NewConsumer(r.client, consumerConfig, messageHandler)

	go func() {
		defer close(messages)
		if err := consumer.Start(ctx); err != nil {
			if err != context.Canceled && err != context.DeadlineExceeded {
				zlog.Logger.Error().Msg("consumer error: " + err.Error())
			}
		}
	}()

	return messages, nil
}

func (r *RabbitMq) Close() error {
	return r.client.Close()
}
