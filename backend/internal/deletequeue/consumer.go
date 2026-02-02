package deletequeue

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConsumer implements DeleteJobConsumer using RabbitMQ.
type RabbitMQConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

// NewRabbitMQConsumer creates a new RabbitMQ consumer.
// Connects to RabbitMQ at the given URL and declares the queue.
func NewRabbitMQConsumer(amqpURL, queueName string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare queue (idempotent: creates if not exists)
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Set QoS: prefetch 1 message at a time for fair distribution
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

// Start subscribes to the delete queue and invokes the handler for each message.
// Only acks when handler returns nil; nacks on error (triggers retry/DLQ per broker policy).
// Runs until ctx is cancelled.
func (c *RabbitMQConsumer) Start(ctx context.Context, handler func(context.Context, DeleteTaskMessage) error) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer tag (empty = auto-generated)
		false,       // auto-ack (false = manual ack)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	log.Printf("[deletequeue] RabbitMQ consumer started for queue: %s", c.queueName)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[deletequeue] Consumer context cancelled, stopping")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("[deletequeue] Message channel closed")
				return nil
			}

			// Deserialize message
			var deleteMsg DeleteTaskMessage
			if err := json.Unmarshal(msg.Body, &deleteMsg); err != nil {
				log.Printf("[Consumer] Failed to unmarshal message: %v", err)
				msg.Nack(false, false) // reject, don't requeue (malformed message)
				continue
			}

			// Process message
			if err := handler(ctx, deleteMsg); err != nil {
				log.Printf("[Consumer] Handler error for task %s: %v (will retry)", deleteMsg.TaskUUID, err)
				// Nack with requeue=true to retry
				msg.Nack(false, true)
				continue
			}

			// Success: ack the message
			msg.Ack(false)
			log.Printf("[Consumer] Successfully processed delete job for task %s", deleteMsg.TaskUUID)
		}
	}
}

// Close closes the RabbitMQ connection and channel.
func (c *RabbitMQConsumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
