package mq

import (
	"golang.org/x/net/context"
)

// Consumer is representation of a message queue consumer.
type Consumer interface {
	// Name return the name of the consumer
	Name() string
	// Topic return which topic is the consumer subscribed to
	Topic() string
	// Group return which group is the consumer belonging to
	Group() string
	// Size return the maximum size of the consumer could process
	Size() int64
	// Consume consumes the message from the publisher
	Consume(ctx context.Context, id string, value any) error
}

// Queue define a set of methods that message queue handler should implement
type Queue interface {
	// Subscribe register consumer itself into Queue then it could receive messages from the specified topic and group
	Subscribe(consumer Consumer) error
	// Publish publishes a message into the specified topic.
	// maxLen is the maximum size of the queue could contain, so add a new entry but will also evict old entries if queue is full,
	// there is no limit if it is zero.
	Publish(ctx context.Context, topic string, value any, maxLen int64) (id string, err error)
	// Start message listening for queue
	Start(ctx context.Context)
	// Close closed the listening
	Close() error
}
