package interfaces

import "context"

// Handler processes a single subscribed message payload.
type Handler func(ctx context.Context, payload []byte) error

// Subscriber registers interest in a topic.
type Subscriber interface {
	Subscribe(topic string, handler Handler) error
}
