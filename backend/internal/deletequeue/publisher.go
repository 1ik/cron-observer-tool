package deletequeue

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher implements DeleteJobPublisher using RabbitMQ.
type RabbitMQPublisher struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher.
// Connects to RabbitMQ at the given URL and declares the queue.
func NewRabbitMQPublisher(amqpURL, queueName string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare queue (idempotent: creates if not exists, same as consumer)
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

	return &RabbitMQPublisher{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

// PublishDeleteTask serializes the message to JSON and publishes it to the delete job queue.
// Returns an error if serialization or publishing fails.
func (p *RabbitMQPublisher) PublishDeleteTask(ctx context.Context, msg DeleteTaskMessage) error {
	// Serialize message to JSON
	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[deletequeue] Failed to marshal DeleteTaskMessage: %v", err)
		return err
	}

	// Publish to queue
	err = p.channel.PublishWithContext(
		ctx,
		"",          // exchange (empty = default/direct exchange)
		p.queueName, // routing key (queue name)
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			// Why persistent for delete jobs?
			// For durable task deletion, we want persistent messages because:
			// Delete jobs are critical: if lost, tasks may remain in PENDING_DELETE indefinitely
			// Reliability: survives RabbitMQ restarts
			// Consistency: matches the durable queue (durable: true)
		},
	)
	if err != nil {
		log.Printf("[deletequeue] Failed to publish delete job for task %s: %v", msg.TaskUUID, err)
		return err
	}

	log.Printf("[deletequeue] Published delete job for task %s to queue %s", msg.TaskUUID, p.queueName)
	return nil
}

// Close closes the RabbitMQ connection and channel.
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
