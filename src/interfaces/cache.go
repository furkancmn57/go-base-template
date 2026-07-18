package interfaces

import (
	"context"
	"time"
)

// Cache is the port for read-through / write caching.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
