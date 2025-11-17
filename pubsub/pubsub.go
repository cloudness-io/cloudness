package pubsub

import "context"

type Publisher interface {
	// Publish message to topic
	Publish(ctx context.Context, topic string, payload []byte, options ...PublishOption) error
}

type Subscriber interface {
	// Subscribe consumer to process the topic with payload, this should be blocking operation.
	Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, options ...SubscribeOption) Consumer
}

type PubSub interface {
	Publisher
	Subscriber
}

type Consumer interface {
	// Subscribe(ctx context.Context, topics ...string) error
	// Publish(ctx context.Context, topics ...string) error
	Close() error
}
