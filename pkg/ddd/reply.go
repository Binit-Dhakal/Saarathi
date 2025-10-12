package ddd

import (
	"time"

	"github.com/google/uuid"
)

type ReplyPayload any

type Reply interface {
	IDer
	ReplyName() string
	Payload() ReplyPayload
	OccuredAt() time.Time
}

type reply struct {
	Entity
	payload   ReplyPayload
	occuredAt time.Time
}

var _ Reply = (*reply)(nil)

func NewReply(name string, payload ReplyPayload) Reply {
	reply := &reply{
		Entity:    NewEntity(uuid.NewString(), name),
		payload:   payload,
		occuredAt: time.Now(),
	}
	return reply
}

func (r reply) ID() string            { return r.Entity.ID() }
func (r reply) ReplyName() string     { return r.Entity.EntityName() }
func (r reply) Payload() ReplyPayload { return r.payload }
func (r reply) OccuredAt() time.Time  { return r.occuredAt }
