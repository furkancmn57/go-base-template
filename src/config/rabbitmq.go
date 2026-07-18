package config

// RabbitMQ holds RabbitMQ connection settings.
type RabbitMQ struct {
	URL string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
}
