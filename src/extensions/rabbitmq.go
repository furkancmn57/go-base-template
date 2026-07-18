package extensions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/furkancmn57/go-base-template/src/config"
	"github.com/furkancmn57/go-base-template/src/constants"
	"github.com/furkancmn57/go-base-template/src/interfaces"
)

// RabbitMQ implements interfaces.Publisher and interfaces.Subscriber.
type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

// AddRabbitMQ dials RabbitMQ and declares the topic exchange.
func AddRabbitMQ(cfg config.RabbitMQ) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: failed to dial: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq: failed to open channel: %w", err)
	}

	if err := channel.ExchangeDeclare(
		constants.RabbitMQExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq: failed to declare exchange: %w", err)
	}

	return &RabbitMQ{
		conn:     conn,
		channel:  channel,
		exchange: constants.RabbitMQExchange,
	}, nil
}

// Publish marshals payload as JSON and sends it (routing key = topic).
func (r *RabbitMQ) Publish(ctx context.Context, topic string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rabbitmq: failed to marshal payload for topic %q: %w", topic, err)
	}

	if err := r.channel.PublishWithContext(
		ctx,
		r.exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		return fmt.Errorf("rabbitmq: failed to publish topic %q: %w", topic, err)
	}
	return nil
}

// Subscribe declares a durable queue bound to the topic and consumes in the background.
func (r *RabbitMQ) Subscribe(topic string, handler interfaces.Handler) error {
	queueName := fmt.Sprintf("%s.queue", topic)

	queue, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("rabbitmq: failed to declare queue %q: %w", queueName, err)
	}

	if err := r.channel.QueueBind(queue.Name, topic, r.exchange, false, nil); err != nil {
		return fmt.Errorf("rabbitmq: failed to bind queue %q to topic %q: %w", queue.Name, topic, err)
	}

	messages, err := r.channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("rabbitmq: failed to start consuming queue %q: %w", queue.Name, err)
	}

	go func() {
		for delivery := range messages {
			if err := handler(context.Background(), delivery.Body); err != nil {
				log.Printf("rabbitmq: handler error for topic %q: %v", topic, err)
				_ = delivery.Nack(false, false)
				continue
			}
			_ = delivery.Ack(false)
		}
	}()

	return nil
}

// Close shuts down the channel and connection.
func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		_ = r.conn.Close()
		return fmt.Errorf("rabbitmq: failed to close channel: %w", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("rabbitmq: failed to close connection: %w", err)
	}
	return nil
}
