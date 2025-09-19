package am

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RequestReplyBus interface {
	Request(ctx context.Context, topic string, cmd ddd.Command) (ddd.Reply, error)
	Reply(ctx context.Context, topic string, handler CommandMessageHandler, options ...SubscriberOption) error
}

type requestReplyBus struct {
	broker RequestTransport
	reg    registry.Registry
}

var _ RequestReplyBus = (*requestReplyBus)(nil)

func NewRequestReplyBus(t RequestTransport, reg registry.Registry) RequestReplyBus {
	return &requestReplyBus{
		broker: t,
		reg:    reg,
	}
}

func (b *requestReplyBus) Request(ctx context.Context, topic string, cmd ddd.Command) (ddd.Reply, error) {
	emptyReply := ddd.NewReply("", "", "")
	payload, err := b.reg.Serialize(cmd.CommandName(), cmd.Payload())
	if err != nil {
		return emptyReply, err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:   payload,
		OccuredAt: timestamppb.New(cmd.OccuredAt()),
	})

	if err != nil {
		return emptyReply, err
	}

	msg, err := b.broker.Request(ctx, topic, rawMessage{
		id:   cmd.ID(),
		name: cmd.CommandName(),
		data: data,
	})

	if err != nil {
		return emptyReply, err
	}

	p, err := b.reg.Deserialize(msg.MessageName(), msg.Data())
	if err != nil {
		return emptyReply, err
	}

	return ddd.NewReply(msg.ID(), msg.MessageName(), p), nil
}

func (b *requestReplyBus) Reply(ctx context.Context, topic string, handler CommandMessageHandler, options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters := make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[IncomingRawMessage](func(ctx context.Context, msg IncomingRawMessage) error {
		var commandData CommandMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &commandData)
		if err != nil {
			return err
		}

		payload, err := b.reg.Deserialize(msg.MessageName(), commandData.GetPayload())
		if err != nil {
			return err
		}

		cmd := &commandMessage{
			id:         msg.ID(),
			name:       msg.MessageName(),
			payload:    payload,
			occurredAt: commandData.GetOccuredAt().AsTime(),
			msg:        msg,
		}

		reply, err := handler.HandleMessage(ctx, cmd)
		if err != nil {
			return err
		}

		data, err := b.reg.Serialize(reply.ReplyName(), reply.Payload())
		if err != nil {
			return err
		}

		_ = b.broker.Publish(ctx, msg.ID(), rawMessage{
			id:   reply.ID(),
			name: reply.ReplyName(),
			data: data,
		})

		if cfg.AckType() == AckTypeManual {
			_ = msg.Ack()
		}

		return nil
	})

	return b.broker.Subscribe(ctx, topic, fn, options...)
}
