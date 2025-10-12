package am

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

const (
	FailureReply = "am.Failure"
	SuccessReply = "am.Success"
)

type IncomingReplyMessage interface {
	IncomingMessage
	ddd.Reply
}

type ReplyPublisher interface {
	Publish(ctx context.Context, topicName string, msg ddd.Reply) error
}

type ReplySubscriber interface {
	Subscribe(ctx context.Context, topicName string, handler MessageHandler[IncomingReplyMessage], options ...SubscriberOption) error
	Unsubscribe() error
}

type replyMessage struct {
	id         string
	name       string
	payload    ddd.ReplyPayload
	occurredAt time.Time
	msg        IncomingMessage
}

func (r replyMessage) ID() string                { return r.id }
func (r replyMessage) ReplyName() string         { return r.name }
func (r replyMessage) Payload() ddd.ReplyPayload { return r.payload }
func (r replyMessage) OccuredAt() time.Time      { return r.occurredAt }
func (r replyMessage) MessageName() string       { return r.msg.MessageName() }
func (r replyMessage) Ack() error                { return r.msg.Ack() }
func (r replyMessage) NAck() error               { return r.msg.NAck() }
func (r replyMessage) Extend() error             { return r.msg.Extend() }
func (r replyMessage) Kill() error               { return r.msg.Kill() }
