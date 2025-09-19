package am

import "context"

type Transport interface {
	Publish(ctx context.Context, topic string, msg RawMessage) error
	Subscribe(ctx context.Context, topic string, handler RawMessageHandler, options ...SubscriberOption) error
}

// for native nats request-reply
type RequestTransport interface {
	Transport
	Request(ctx context.Context, topic string, msg RawMessage) (RawMessage, error)
}
