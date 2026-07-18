package extensions

import (
	goredis "github.com/redis/go-redis/v9"

	"github.com/furkancmn57/go-base-template/src/config"
	"github.com/furkancmn57/go-base-template/src/services/cache"
)

// AddRedis opens and pings a Redis client.
func AddRedis(cfg config.Redis) (*goredis.Client, error) {
	return cache.NewClient(cfg)
}
