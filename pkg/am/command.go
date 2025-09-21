package am

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type CommandMessage interface {
	Message
	ddd.Command
}

type IncomingCommandMessage interface {
	IncomingMessage
	ddd.Command
}

type CommandMessageHandler interface {
	HandleMessage(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error)
}

type CommandMessageHandlerFunc func(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error)

type CommandSender interface {
	Send(ctx context.Context, topicName string, cmd ddd.Command) (RawMessage, error)
}

type CommandSubscriber interface {
	Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error
	Unsubscribe() error
}

type commandMessage struct {
	id         string
	name       string
	payload    ddd.CommandPayload
	occurredAt time.Time
	msg        IncomingMessage
}

func (f CommandMessageHandlerFunc) HandleMessage(ctx context.Context, msg IncomingCommandMessage) (ddd.Reply, error) {
	return f(ctx, msg)
}

var _ CommandMessage = (*commandMessage)(nil)

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Payload() ddd.CommandPayload { return c.payload }
func (c commandMessage) OccuredAt() time.Time        { return c.occurredAt }
func (c commandMessage) MessageName() string         { return c.msg.MessageName() }
func (c commandMessage) Ack() error                  { return c.msg.Ack() }
func (c commandMessage) NAck() error                 { return c.msg.NAck() }
func (c commandMessage) Extend() error               { return c.msg.Extend() }
func (c commandMessage) Kill() error                 { return c.msg.Kill() }
