package messagebus

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
)

type Publisher interface {
	Publish(ctx context.Context, exchange string, routingKey string, message events.Event) error
}

type Consumer interface {
	Consume(ctx context.Context, queue string)
}

type EventHandler func(ctx context.Context, event events.Event) error
