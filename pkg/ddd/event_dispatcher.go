package ddd

import "context"

type EventHandler[T Event] interface {
	HandleEvent(ctx context.Context, event T) error
}

type EventSubscriber[T Event] interface {
	Subscribe(handler EventHandler[T], events ...string)
}

type EventPublisher[T Event] interface {
	Publish(ctx context.Context, events ...T) error
}

