package deletequeue

import "context"

// DeleteJobPublisher is a broker-agnostic interface for publishing delete job messages.
// Implementations may use RabbitMQ, SQS, Redis, or any other message broker;
// the rest of the code stays independent of the specific broker.
type DeleteJobPublisher interface {
	PublishDeleteTask(ctx context.Context, msg DeleteTaskMessage) error
}

// DeleteJobConsumer is a broker-agnostic interface for consuming delete job messages.
// Start subscribes to the delete queue and invokes the handler for each message.
// The handler should return nil to acknowledge the message, or an error to trigger retry/DLQ per broker policy.
type DeleteJobConsumer interface {
	Start(ctx context.Context, handler func(context.Context, DeleteTaskMessage) error) error
}
