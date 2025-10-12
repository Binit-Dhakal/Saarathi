package am

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type Message interface {
	ddd.IDer
	MessageName() string
}

type IncomingMessage interface {
	Message
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}

type MessageHandler[I IncomingMessage] interface {
	HandleMessage(ctx context.Context, msg I) error
}

type MessageHandlerFunc[I IncomingMessage] func(ctx context.Context, msg I) error

func (f MessageHandlerFunc[I]) HandleMessage(ctx context.Context, msg I) error {
	return f(ctx, msg)
}
