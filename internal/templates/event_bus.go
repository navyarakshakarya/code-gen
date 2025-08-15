package templates

// RabbitMQ Event Bus Template
const EventBusTemplate = `package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventBus defines the interface for an event bus
type EventBus interface {
	Publish(exchange, routingKey string, event interface{}) error
	Subscribe(ctx context.Context, exchange, queue, routingKey string, handler func([]byte) error) error
	Close() error
}

// RabbitMQBus implements EventBus using RabbitMQ
type RabbitMQBus struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQBus creates a new RabbitMQ event bus
func NewRabbitMQBus(url string) (*RabbitMQBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQBus{conn: conn, channel: ch}, nil
}

func (r *RabbitMQBus) declareExchange(exchange string) error {
	return r.channel.ExchangeDeclare(
		exchange,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	)
}

// Publish publishes an event to the specified exchange
func (r *RabbitMQBus) Publish(exchange, routingKey string, event interface{}) error {
	if err := r.declareExchange(exchange); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	return r.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

// Subscribe subscribes to events and processes them until ctx is cancelled
func (r *RabbitMQBus) Subscribe(ctx context.Context, exchange, queue, routingKey string, handler func([]byte) error) error {
	if err := r.declareExchange(exchange); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	q, err := r.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	if err := r.channel.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return // channel closed
				}
				if err := handler(msg.Body); err != nil {
					log.Printf("error handling message: %v", err)
					_ = msg.Nack(false, true) // requeue
				} else {
					_ = msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQBus) Close() error {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// Event represents a domain event
type Event struct {
	ID        string                ` + "`json:\"id\"`" + `
	Type      string                ` + "`json:\"type\"`" + `
	Timestamp int64                 ` + "`json:\"timestamp\"`" + `
	Data      map[string]interface{} ` + "`json:\"data\"`" + `
}

`
