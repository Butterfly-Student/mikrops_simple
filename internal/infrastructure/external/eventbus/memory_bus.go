package eventbus

import (
	"context"
	"github.com/alijayanet/gembok-backend/internal/domain/events"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
	"sync"
)

type MemoryEventBus struct {
	mu          sync.RWMutex
	handlers    map[string][]EventHandler
	allHandlers []EventHandler
}

func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers:    make(map[string][]EventHandler),
		allHandlers: nil,
	}
}

func (b *MemoryEventBus) Publish(ctx context.Context, event events.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, exists := b.handlers[event.GetType()]; exists {
		for _, handler := range handlers {
			go func(h EventHandler) {
				if err := h(ctx, event); err != nil {
					logger.Error("Event handler error",
						zap.String("event_type", event.GetType()),
						zap.String("event_id", event.GetID()),
						zap.Error(err))
				}
			}(handler)
		}
	}

	for _, handler := range b.allHandlers {
		go func(h EventHandler) {
			if err := h(ctx, event); err != nil {
				logger.Error("All-event handler error",
					zap.String("event_type", event.GetType()),
					zap.String("event_id", event.GetID()),
					zap.Error(err))
			}
		}(handler)
	}

	logger.Debug("Event published",
		zap.String("event_type", event.GetType()),
		zap.String("event_id", event.GetID()))

	return nil
}

func (b *MemoryEventBus) Subscribe(ctx context.Context, eventType string, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)

	logger.Debug("Event handler subscribed",
		zap.String("event_type", eventType))

	return nil
}

func (b *MemoryEventBus) SubscribeAll(ctx context.Context, handler EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.allHandlers = append(b.allHandlers, handler)

	logger.Debug("All-event handler subscribed")
	return nil
}

func (b *MemoryEventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = make(map[string][]EventHandler)
	b.allHandlers = nil

	logger.Info("Memory event bus closed")
	return nil
}
