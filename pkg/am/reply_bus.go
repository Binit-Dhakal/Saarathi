package am

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReplyBus interface {
	ReplyPublisher
	ReplySubscriber
}

type replyBus struct {
	reg    registry.Registry
	broker Transport
}

var _ ReplyBus = (*replyBus)(nil)

func NewReplyBus(reg registry.Registry, broker Transport) *replyBus {
	return &replyBus{
		reg:    reg,
		broker: broker,
	}
}

func (r *replyBus) Publish(ctx context.Context, topic string, msg ddd.Reply) (err error) {
	var payload []byte

	if msg.ReplyName() != SuccessReply && msg.ReplyName() != FailureReply {
		payload, err = r.reg.Serialize(msg.ReplyName(), msg.Payload())
		if err != nil {
			return err
		}
	}

	data, err := proto.Marshal(&ReplyMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(msg.OccuredAt()),
	})
	if err != nil {
		return err
	}

	return r.broker.Publish(ctx, topic, &rawMessage{
		id:   msg.ID(),
		name: msg.ReplyName(),
		data: data,
	})
}

func (r *replyBus) Subscribe(ctx context.Context, topicName string, handler MessageHandler[IncomingReplyMessage], options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	replyHandler := MessageHandlerFunc[IncomingRawMessage](func(ctx context.Context, msg IncomingRawMessage) error {
		var replyData ReplyMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &replyData)
		if err != nil {
			return err
		}

		replyName := msg.MessageName()

		payload, err := r.reg.Deserialize(replyName, replyData.GetPayload())
		if err != nil {
			return err
		}

		replyMsg := replyMessage{
			id:         msg.ID(),
			name:       replyName,
			payload:    payload,
			occurredAt: replyData.OccurredAt.AsTime(),
			msg:        msg,
		}

		return handler.HandleMessage(ctx, replyMsg)
	})

	return r.broker.Subscribe(ctx, topicName, replyHandler, options...)
}

func (r *replyBus) Unsubscribe() error {
	return r.broker.Unsubscribe()
}
