package command

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"google.golang.org/protobuf/proto"
)

type Requestor interface {
	Request(ctx context.Context, subject string, cmd Command) (Reply, error)
}

type Responder interface {
	Respond(subject string, handler func(ctx context.Context, cmd Command)) (Reply, error)
}

type CommandPublisher struct {
	broker Broker
	reg    registry.Registry
}

var _ Requestor = (*CommandPublisher)(nil)

func NewCommandPublisher(reg registry.Registry, broker Broker) *CommandPublisher {
	return &CommandPublisher{reg: reg, broker: broker}
}

func (p *CommandPublisher) Request(ctx context.Context, subject string, cmd Command) (Reply, error) {
	payload, err := p.reg.Serialize(cmd.CommandName(), cmd.Payload())
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload: payload,
	})
	if err != nil {
		return nil, err
	}

	p.broker.Request(ctx, subject, am.NewRawMessage(cmd.ID(), cmd.CommandName(), data))

}

type CommandSubscriber struct {
	reg     registry.Registry
	stream  RawMessageStream
	subject string
}

func NewCommandSubscriber(reg registry.Registry, stream RawMessageStream, subject string) *CommandSubscriber {
	return &CommandSubscriber{reg: reg, stream: stream, subject: subject}
}
