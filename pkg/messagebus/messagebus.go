package messagebus

import "context"

type Publisher interface {
	Publish(ctx context.Context, exchange string, topic string, message any) error
}

type Consumer interface {
	Consume(ctx context.Context, queue string)
}
