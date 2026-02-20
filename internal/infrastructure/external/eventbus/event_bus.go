package eventbus

import (
	"context"
	"github.com/alijayanet/gembok-backend/internal/domain/events"
)

type EventBus interface {
	Publish(ctx context.Context, event events.Event) error
	Subscribe(ctx context.Context, eventType string, handler EventHandler) error
	SubscribeAll(ctx context.Context, handler EventHandler) error
	Close() error
}

type EventHandler func(ctx context.Context, event events.Event) error

type EventConfig struct {
	Type     string
	RabbitMQ *RabbitMQConfig
}

type RabbitMQConfig struct {
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
	Durable    bool
}
