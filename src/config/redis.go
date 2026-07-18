package config

import "fmt"

// Redis holds Redis connection settings.
type Redis struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

// Addr returns host:port for the redis client.
func (c Redis) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
