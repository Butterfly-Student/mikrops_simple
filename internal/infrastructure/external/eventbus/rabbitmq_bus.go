package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alijayanet/gembok-backend/internal/domain/events"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQEventBus struct {
	conn       *amqp091.Connection
	channel    *amqp091.Channel
	exchange   string
	queue      string
	routingKey string
}

func NewRabbitMQEventBus(config *RabbitMQConfig) (*RabbitMQEventBus, error) {
	conn, err := amqp091.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = channel.ExchangeDeclare(
		config.Exchange,
		"topic",
		config.Durable,
		true,
		false,
		false,
		amqp091.Table{},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		config.Queue,
		config.Durable,
		false,
		false,
		false,
		amqp091.Table{},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	err = channel.QueueBind(
		queue.Name,
		config.RoutingKey,
		config.Exchange,
		false,
		amqp091.Table{},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	bus := &RabbitMQEventBus{
		conn:       conn,
		channel:    channel,
		exchange:   config.Exchange,
		queue:      queue.Name,
		routingKey: config.RoutingKey,
	}

	logger.Info("RabbitMQ event bus connected",
		zap.String("exchange", config.Exchange),
		zap.String("queue", config.Queue))

	return bus, nil
}

func (b *RabbitMQEventBus) Publish(ctx context.Context, event events.Event) error {
	eventJSON, err := event.ToJSON()
	if err != nil {
		return err
	}

	err = b.channel.PublishWithContext(
		ctx,
		b.exchange,
		b.routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         eventJSON,
			DeliveryMode: amqp091.Persistent,
		},
	)

	if err != nil {
		logger.Error("Failed to publish event to RabbitMQ",
			zap.String("event_type", event.GetType()),
			zap.Error(err))
		return err
	}

	logger.Debug("Event published to RabbitMQ",
		zap.String("event_type", event.GetType()),
		zap.String("event_id", event.GetID()))

	return nil
}

func (b *RabbitMQEventBus) Subscribe(ctx context.Context, eventType string, handler EventHandler) error {
	msgs, err := b.channel.Consume(
		b.queue,
		"",
		false,
		false,
		false,
		false,
		amqp091.Table{},
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var baseEvent events.BaseEvent
			if err := json.Unmarshal(msg.Body, &baseEvent); err != nil {
				logger.Error("Failed to unmarshal event",
					zap.Error(err))
				msg.Nack(false, true)
				continue
			}

			if baseEvent.Type != eventType {
				msg.Ack(false)
				continue
			}

			if err := handler(ctx, &baseEvent); err != nil {
				logger.Error("Event handler error",
					zap.String("event_type", baseEvent.Type),
					zap.Error(err))
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	logger.Debug("RabbitMQ event handler subscribed",
		zap.String("event_type", eventType))

	return nil
}

func (b *RabbitMQEventBus) SubscribeAll(ctx context.Context, handler EventHandler) error {
	return errors.New("SubscribeAll not supported in RabbitMQ event bus")
}

func (b *RabbitMQEventBus) Close() error {
	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}

	logger.Info("RabbitMQ event bus closed")
	return nil
}
