package am

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommandBus interface {
	CommandSender
	CommandSubscriber
}

type commandBus struct {
	reg    registry.Registry
	broker RequestTransport
}

var _ CommandBus = (*commandBus)(nil)

func NewCommandBus(reg registry.Registry, broker RequestTransport) CommandBus {
	return &commandBus{
		reg:    reg,
		broker: broker,
	}
}

func (b *commandBus) Send(ctx context.Context, topicName string, cmd ddd.Command) (RawMessage, error) {
	payload, err := b.reg.Serialize(cmd.CommandName(), cmd.Payload())
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:   payload,
		OccuredAt: timestamppb.New(cmd.OccuredAt()),
	})

	if err != nil {
		return nil, err
	}

	return b.broker.Request(ctx, topicName, &rawMessage{
		id:   cmd.ID(),
		name: cmd.CommandName(),
		data: data,
	})
}

func (b *commandBus) Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters := make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	replyHandler := func(ctx context.Context, req RawMessage) (RawMessage, error) {
		var commandData CommandMessageData

		if filters != nil {
			if _, exists := filters[req.MessageName()]; !exists {
				return nil, nil
			}
		}

		err := proto.Unmarshal(req.Data(), &commandData)
		if err != nil {
			return nil, err
		}

		commandName := req.MessageName()

		payload, err := b.reg.Deserialize(commandName, commandData.GetPayload())
		if err != nil {
			return nil, err
		}

		commandMsg := commandMessage{
			id:         req.ID(),
			name:       commandName,
			payload:    payload,
			occurredAt: commandData.GetOccuredAt().AsTime(),
			msg:        nil,
		}

		resp, err := handler.HandleMessage(ctx, commandMsg)
		if err != nil {
			return nil, err
		}

		respPayload, err := b.reg.Serialize(resp.ReplyName(), resp.Payload())
		if err != nil {
			return nil, err
		}

		replyProto := ReplyMessageData{
			Payload:    respPayload,
			OccurredAt: timestamppb.Now(),
		}

		bts, err := proto.Marshal(&replyProto)
		if err != nil {
			return nil, err
		}

		return rawMessage{
			id:   resp.ID(),
			name: resp.ReplyName(),
			data: bts,
		}, nil
	}

	return b.broker.Reply(topicName, replyHandler, options...)
}

func (b *commandBus) Unsubscribe() error {
	return b.broker.Close()
}
