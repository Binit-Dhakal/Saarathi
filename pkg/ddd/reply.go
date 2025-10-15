package ddd

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ReplyHandler[T Reply] interface {
	HandleReply(ctx context.Context, reply T) error
}

type ReplyHandlerFunc[T Reply] func(ctx context.Context, reply T) error

type ReplyOption interface {
	configureReply(*reply)
}

type ReplyPayload any

type Reply interface {
	IDer
	ReplyName() string
	Payload() ReplyPayload
	Metadata() Metadata
	OccuredAt() time.Time
}

type reply struct {
	Entity
	payload   ReplyPayload
	metadata  Metadata
	occuredAt time.Time
}

var _ Reply = (*reply)(nil)

func NewReply(name string, payload ReplyPayload, options ...ReplyOption) Reply {
	reply := &reply{
		Entity:    NewEntity(uuid.NewString(), name),
		payload:   payload,
		occuredAt: time.Now(),
	}

	for _, option := range options {
		option.configureReply(reply)
	}

	return reply
}

func (r reply) ID() string            { return r.Entity.ID() }
func (r reply) ReplyName() string     { return r.Entity.EntityName() }
func (r reply) Payload() ReplyPayload { return r.payload }
func (r reply) OccuredAt() time.Time  { return r.occuredAt }
func (r reply) Metadata() Metadata    { return r.metadata }

func (f ReplyHandlerFunc[T]) HandleReply(ctx context.Context, reply T) error {
	return f(ctx, reply)
}
