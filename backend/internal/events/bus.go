package events

import (
	"sync"
)

// EventBus manages event subscriptions and publishing
type EventBus struct {
	subscribers map[EventType][]chan Event
	mu          sync.RWMutex
	bufferSize  int
}

// NewEventBus creates a new EventBus with the specified buffer size for channels
func NewEventBus(bufferSize int) *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]chan Event),
		bufferSize:  bufferSize,
	}
}

// Subscribe creates a subscription channel for a specific event type
func (b *EventBus) Subscribe(eventType EventType) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, b.bufferSize)
	b.subscribers[eventType] = append(b.subscribers[eventType], ch)

	return ch
}

// Publish sends an event to all subscribers of that event type
func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	subscribers := b.subscribers[event.Type]
	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
			// Channel is full, skip to avoid blocking
			// In production, you might want to log this
		}
	}
}

// Close closes all subscriber channels
func (b *EventBus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, channels := range b.subscribers {
		for _, ch := range channels {
			close(ch)
		}
	}
	b.subscribers = make(map[EventType][]chan Event)
}
