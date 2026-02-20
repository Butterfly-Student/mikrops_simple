package eventbus

import (
	"fmt"
)

func NewEventBus(config *EventConfig) (EventBus, error) {
	switch config.Type {
	case "rabbitmq":
		if config.RabbitMQ == nil {
			return nil, fmt.Errorf("rabbitmq config is required when type is 'rabbitmq'")
		}
		return NewRabbitMQEventBus(config.RabbitMQ)
	case "memory":
		fallthrough
	default:
		return NewMemoryEventBus(), nil
	}
}
