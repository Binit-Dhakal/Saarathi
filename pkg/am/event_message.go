package am

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type EventMessage interface {
	MessageBase
	ddd.Event
}

type IncomingEventMessage interface {
	IncomingMessageBase
	ddd.Event
}

type EventPublisher interface {
	Publish(ctx context.Context, topicName string, event ddd.Event) error
}

type eventPublisher struct {
	publisher MessagePublisher
}

type eventMessage struct {
	id        string
	name      string
	payload   ddd.EventPayload
	occuredAt time.Time
	msg       IncomingMessageBase
}

var _ EventMessage = (*eventMessage)(nil)
var _ EventPublisher = (*eventPublisher)(nil)

func NewEventPublisher(msgPublisher MessagePublisher) EventPublisher {
	return eventPublisher{
		publisher: msgPublisher,
	}
}

func (s eventPublisher) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, topicName, message{
		id:       event.ID(),
		name:     event.EventName(),
		subject:  topicName,
		data:     data,
		metadata: event.Metadata(),
		sentAt:   time.Now(),
	})
}

func (e eventMessage) ID() string                { return e.id }
func (e eventMessage) EventName() string         { return e.name }
func (e eventMessage) Payload() ddd.EventPayload { return e.payload }
func (e eventMessage) Metadata() ddd.Metadata    { return e.msg.Metadata() }
func (e eventMessage) OccuredAt() time.Time      { return e.occuredAt }
func (e eventMessage) Subject() string           { return e.msg.Subject() }
func (e eventMessage) MessageName() string       { return e.msg.MessageName() }
func (e eventMessage) SentAt() time.Time         { return e.msg.SentAt() }
func (e eventMessage) ReceivedAt() time.Time     { return e.msg.ReceivedAt() }
func (e eventMessage) Ack() error                { return e.msg.Ack() }
func (e eventMessage) NAck() error               { return e.msg.NAck() }
func (e eventMessage) Extend() error             { return e.msg.Extend() }
func (e eventMessage) Kill() error               { return e.msg.Kill() }
