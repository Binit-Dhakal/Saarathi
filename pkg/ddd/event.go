package ddd

import (
	"time"

	"github.com/google/uuid"
)

type EventOption interface {
	configureEvent(*event)
}

type Event interface {
	IDer
	EventName() string
	Payload() EventPayload
	Metadata() Metadata
	OccuredAt() time.Time
}

type EventPayload any

type event struct {
	Entity
	payload   EventPayload
	metadata  Metadata
	occuredAt time.Time
}

func NewEvent(name string, payload EventPayload, options ...EventOption) Event {
	return newEvent(name, payload, options...)
}

func newEvent(name string, payload EventPayload, options ...EventOption) event {
	e := event{
		Entity:    NewEntity(uuid.NewString(), name),
		payload:   payload,
		metadata:  make(Metadata),
		occuredAt: time.Now(),
	}

	for _, option := range options {
		option.configureEvent(&e)
	}

	return e
}

func (e event) EventName() string     { return e.Entity.EntityName() }
func (e event) ID() string            { return e.Entity.ID() }
func (e event) Payload() EventPayload { return e.payload }
func (e event) Metadata() Metadata    { return e.metadata }
func (e event) OccuredAt() time.Time  { return e.occuredAt }
