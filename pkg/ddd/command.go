package ddd

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CommandPayload any

type CommandOption interface {
	configureCommand(*command)
}

type Command interface {
	IDer
	CommandName() string
	Payload() CommandPayload
	Metadata() Metadata
	OccuredAt() time.Time
}

type CommandHandler interface {
	HandleCommand(ctx context.Context, cmd Command) (Reply, error)
}

type CommandHandlerFunc func(ctx context.Context, cmd Command) (Reply, error)

type command struct {
	Entity
	payload   CommandPayload
	metadata  Metadata
	occuredAt time.Time
}

var _ Command = (*command)(nil)

func NewCommand(name string, payload CommandPayload, options ...CommandOption) Command {
	command := &command{
		Entity:    NewEntity(uuid.NewString(), name),
		payload:   payload,
		occuredAt: time.Now(),
	}

	for _, option := range options {
		option.configureCommand(command)
	}

	return command
}

func (c command) ID() string              { return c.Entity.ID() }
func (c command) CommandName() string     { return c.Entity.EntityName() }
func (c command) Payload() CommandPayload { return c.payload }
func (c command) OccuredAt() time.Time    { return c.occuredAt }
func (c command) Metadata() Metadata      { return c.metadata }

func (c CommandHandlerFunc) HandleCommand(ctx context.Context, cmd Command) (Reply, error) {
	return c(ctx, cmd)
}
