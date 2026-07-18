package common

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/furkancmn57/go-base-template/src/constants"
)

// HandleFunc processes a raw messaging payload.
type HandleFunc func(ctx context.Context, payload []byte) error

// BaseConsumer wraps a handler with optional payload logging when
// CONSUMER_LOGS_ENABLED=true.
func BaseConsumer(name string, next HandleFunc) HandleFunc {
	enabled := strings.EqualFold(os.Getenv(constants.ConsumerLogsEnabledEnv), "true")
	return func(ctx context.Context, payload []byte) error {
		if enabled {
			log.Printf("%s: received message (%d bytes)", name, len(payload))
		}
		err := next(ctx, payload)
		if enabled && err == nil {
			log.Printf("%s: message completed", name)
		}
		return err
	}
}
