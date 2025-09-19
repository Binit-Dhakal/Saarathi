package ddd

import (
	"time"

	"github.com/google/uuid"
)

type CommandPayload any

type Command interface {
	IDer
	CommandName() string
	Payload() CommandPayload
	OccuredAt() time.Time
}

type command struct {
	Entity
	payload   CommandPayload
	occuredAt time.Time
}

var _ Command = (*command)(nil)

func NewCommand(name string, payload CommandPayload) Command {
	command := &command{
		Entity:    NewEntity(uuid.NewString(), name),
		payload:   payload,
		occuredAt: time.Now(),
	}
	return command
}

func (c command) ID() string              { return c.Entity.ID() }
func (c command) CommandName() string     { return c.Entity.EntityName() }
func (c command) Payload() CommandPayload { return c.payload }
func (c command) OccuredAt() time.Time    { return c.occuredAt }
