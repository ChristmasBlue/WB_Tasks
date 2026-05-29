package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"task-5/internal/config"
	"task-5/internal/dto"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
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

	// Создаем конфигурацию клиента
	clientConfig := rabbitmq.ClientConfig{
		URL:            url,
		ConnectionName: "task-5",
		ConnectTimeout: 10 * time.Second,
		Heartbeat:      10 * time.Second,
		ReconnectStrat: retry.Strategy{
			Attempts: 0, // бесконечные попытки переподключения
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

	// Создаем клиент
	client, err := rabbitmq.NewClient(clientConfig)
	if err != nil {
		zlog.Logger.Fatal().Msg("could not create rabbitmq client: " + err.Error())
	}

	// Объявляем delayed exchange
	if err := declareDelayedExchange(client); err != nil {
		zlog.Logger.Fatal().Msg("could not declare delayed exchange: " + err.Error())
	}

	// Объявляем очередь и привязываем к exchange
	if err := declareQueue(client); err != nil {
		zlog.Logger.Fatal().Msg("could not declare queue: " + err.Error())
	}

	// Создаем publisher
	publisher := rabbitmq.NewPublisher(client, "bookings", "application/json")

	return &RabbitMq{
		client:    client,
		publisher: publisher,
	}
}

// declareDelayedExchange объявляет exchange с типом x-delayed-message
func declareDelayedExchange(client *rabbitmq.RabbitClient) error {
	ch, err := client.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()

	args := amqp.Table{"x-delayed-type": "direct"}
	return ch.ExchangeDeclare(
		"bookings",
		"x-delayed-message",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		args,
	)
}

// declareQueue объявляет очередь и привязывает её к exchange
func declareQueue(client *rabbitmq.RabbitClient) error {
	ch, err := client.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Объявляем очередь
	_, err = ch.QueueDeclare(
		"bookings",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// Привязываем очередь к exchange
	return ch.QueueBind("bookings", "bookings", "bookings", false, nil)
}

func (r *RabbitMq) Publish(booking dto.QueueMessage) error {
	body, err := json.Marshal(booking)
	if err != nil {
		return fmt.Errorf("could not marshal booking to send to rabbitmq: %w", err)
	}

	ctx := context.Background()

	// Используем опции для добавления заголовков с задержкой
	headers := amqp.Table{
		"x-delay": int64((time.Minute * 15).Milliseconds()),
	}

	return r.publisher.Publish(
		ctx,
		body,
		"bookings",
		rabbitmq.WithHeaders(headers),
	)
}

func (r *RabbitMq) Consume(ctx context.Context) (<-chan []byte, error) {
	messages := make(chan []byte, 100)

	// Создаем message handler
	messageHandler := func(ctx context.Context, delivery amqp.Delivery) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case messages <- delivery.Body:
			return nil
		}
	}

	// Конфигурация consumer
	consumerConfig := rabbitmq.ConsumerConfig{
		Queue:       "bookings",
		ConsumerTag: "task-5-consumer",
		AutoAck:     true, // auto-ack, так как обработка в сервисе
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

	// Создаем consumer
	consumer := rabbitmq.NewConsumer(r.client, consumerConfig, messageHandler)

	// Запускаем consumer в горутине
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

// Close закрывает соединение с RabbitMQ
func (r *RabbitMq) Close() error {
	return r.client.Close()
}
