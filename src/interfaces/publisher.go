package interfaces

import "context"

// Publisher publishes best-effort domain events after a DB commit.
type Publisher interface {
	Publish(ctx context.Context, topic string, payload any) error
}
