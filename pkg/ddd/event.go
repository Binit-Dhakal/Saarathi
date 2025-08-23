package ddd

import (
	"time"

	"github.com/google/uuid"
)

type EventPayload any

type Event interface {
	EventName() string
	Payload() EventPayload
	OccuredAt() time.Time
}

type event struct {
	Entity
	payload   EventPayload
	occuredAt time.Time
}

var _ Event = (*event)(nil)

func NewEvent(name string, payload EventPayload) Event {
	return newEvent(name, payload)
}

func newEvent(name string, payload EventPayload) event {
	evt := event{
		Entity:    NewEntity(uuid.New().String(), name),
		payload:   payload,
		occuredAt: time.Now(),
	}
	return evt
}

func (e event) EventName() string     { return e.EntityName() }
func (e event) Payload() EventPayload { return e.payload }
func (e event) OccuredAt() time.Time  { return e.occuredAt }
