package am

import "context"

type Transport interface {
	Publish(ctx context.Context, topic string, msg RawMessage) error
	Subscribe(ctx context.Context, topic string, handler RawMessageHandler, options ...SubscriberOption) error
	Unsubscribe() error
}

// for native nats request-reply
type RequestTransport interface {
	Request(ctx context.Context, topic string, msg RawMessage) (RawMessage, error)
	Reply(subject string, handler func(ctx context.Context, req RawMessage) (RawMessage, error), options ...SubscriberOption) error
	Close() error
}
