package am

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type MessageBase interface {
	ddd.IDer
	Subject() string
	MessageName() string
	Metadata() ddd.Metadata
	SentAt() time.Time
}

type IncomingMessageBase interface {
	MessageBase
	ReceivedAt() time.Time
	Ack() error
	NAck() error
	Extend() error
	Kill() error
}

type IncomingMessage interface {
	IncomingMessageBase
	Data() []byte
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, msg IncomingMessage) error
}

type MessagePublisher interface {
	Publish(ctx context.Context, topicName string, msg Message) error
}

type MessageSubsciber interface {
	Subscribe(topicName string, handler MessageHandler, options ...SubscriberOption) (Subscription, error)
}

type MessageStream interface {
	MessagePublisher
	MessageSubsciber
}

type Message interface {
	MessageBase
	Data() []byte
}

type message struct {
	id       string
	name     string
	subject  string
	data     []byte
	metadata ddd.Metadata
	sentAt   time.Time
}

type messagePublisher struct {
	publisher MessagePublisher
}

type messageSubscriber struct {
	subscriber MessageSubsciber
}

var _ Message = (*message)(nil)

func (m message) ID() string             { return m.id }
func (m message) Data() []byte           { return m.data }
func (m message) MessageName() string    { return m.name }
func (m message) Metadata() ddd.Metadata { return m.metadata }
func (m message) SentAt() time.Time      { return m.sentAt }
func (m message) Subject() string        { return m.subject }

func NewPublisher(publisher MessagePublisher) MessagePublisher {
	return messagePublisher{
		publisher: publisher,
	}
}

func (m messagePublisher) Publish(ctx context.Context, topicName string, msg Message) error {
	return m.publisher.Publish(ctx, topicName, msg)
}

func NewSubscriber(subscriber MessageSubsciber) messageSubscriber {
	return messageSubscriber{
		subscriber: subscriber,
	}
}

func (m messageSubscriber) Subscribe(topicName string, handler MessageHandler) (Subscription, error) {
	return m.subscriber.Subscribe(topicName, handler)
}
