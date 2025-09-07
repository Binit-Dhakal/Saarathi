package jetstream

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

const maxRetries = 5

type Stream struct {
	streamName string
	js         nats.JetStreamContext
	mu         sync.Mutex
	subs       []*nats.Subscription
	logger     zerolog.Logger
}

var _ am.MessageStream = (*Stream)(nil)

func NewStream(streamName string, js nats.JetStreamContext, logger zerolog.Logger) *Stream {
	return &Stream{
		streamName: streamName,
		js:         js,
		logger:     logger,
	}
}

func (s *Stream) Publish(ctx context.Context, topicName string, msg am.Message) (err error) {
	var data []byte

	data, err = json.Marshal(msg)
	if err != nil {
		return err
	}

	var p nats.PubAckFuture
	p, err = s.js.PublishMsgAsync(&nats.Msg{
		Subject: topicName,
		Data:    data,
	}, nats.MsgId(msg.ID()))

	go func(future nats.PubAckFuture, tries int) {
		var err error
		for {
			select {
			case <-future.Ok():
				return
			case <-future.Err():
				tries = tries - 1
				if tries <= 0 {
					s.logger.Error().Msgf("unable to publish message after %d tries", maxRetries)
					return
				}

				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					s.logger.Error().Err(err).Msg("failed to publish message")
					return
				}
			}
		}
	}(p, maxRetries)

	return
}

func (s *Stream) Subscribe(topicName string, handler am.MessageHandler, options ...am.SubscriberOption) (am.Subscription, error) {
	var err error

	s.mu.Lock()
	defer s.mu.Unlock()

	subCfg := am.NewSubscriberConfig(options)

	opts := []nats.SubOpt{
		nats.MaxDeliver(subCfg.MaxRedeliver()),
	}

	cfg := &nats.ConsumerConfig{
		MaxDeliver:     subCfg.MaxRedeliver(),
		DeliverSubject: topicName,
		FilterSubject:  topicName,
	}

	if groupName := subCfg.GroupName(); groupName != "" {
		cfg.DeliverSubject = groupName
		cfg.DeliverGroup = groupName
		cfg.Durable = groupName

		opts = append(opts, nats.Bind(s.streamName, groupName), nats.Durable(groupName))
	}

	if ackType := subCfg.AckType(); ackType != am.AckTypeAuto {
		ackWait := subCfg.AckWait()

		cfg.AckPolicy = nats.AckExplicitPolicy
		cfg.AckWait = ackWait

		opts = append(opts, nats.AckExplicit(), nats.AckWait(ackWait))
	} else {
		cfg.AckPolicy = nats.AckNonePolicy
		opts = append(opts, nats.AckNone())
	}

	_, err = s.js.AddConsumer(s.streamName, cfg)
	if err != nil {
		return nil, err
	}

	var sub *nats.Subscription

	if groupName := subCfg.GroupName(); groupName == "" {
		sub, err = s.js.Subscribe(topicName, s.handleMsg(subCfg, handler), opts...)
	} else {
		sub, err = s.js.QueueSubscribe(topicName, groupName, s.handleMsg(subCfg, handler), opts...)
	}

	s.subs = append(s.subs, sub)

	return subscription{s: sub}, nil
}

func (s *Stream) Unsubscribe() error {
	for _, sub := range s.subs {
		if !sub.IsValid() {
			continue
		}
		err := sub.Drain()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler) func(*nats.Msg) {
	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	return func(natsMsg *nats.Msg) {
		var m am.Message
		err := json.Unmarshal(natsMsg.Data, &m)
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to unmarshal the *nats.Msg")
			return
		}

		if filters != nil {
			if _, exist := filters[m.MessageName()]; !exist {
				err = natsMsg.Ack()
				if err != nil {
					s.logger.Warn().Err(err).Msg("failed to Ack a filtered message")
				}
				return
			}

		}

		msg := &rawMessage{
			id:         m.ID(),
			name:       m.MessageName(),
			subject:    m.Subject(),
			data:       m.Data(),
			metadata:   m.Metadata(),
			sentAt:     m.SentAt(),
			receivedAt: time.Now(),
			acked:      false,
			ackFn:      func() error { return natsMsg.Ack() },
			nackFn:     func() error { return natsMsg.Nak() },
			extendFn:   func() error { return natsMsg.InProgress() },
			killFn:     func() error { return natsMsg.Term() },
		}

		wCtx, cancel := context.WithTimeout(context.Background(), cfg.AckWait())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err = msg.Ack()
			if err != nil {
				s.logger.Warn().Err(err).Msg("failed to auto-ack a message")
			}
		}

		select {
		case err = <-errc:
			if err != nil {
				if ackErr := msg.Ack(); ackErr != nil {
					s.logger.Warn().Err(err).Msg("failed to ack a message")
				}
				return
			}

			s.logger.Error().Err(err).Msg("error while handling message")
			if nakErr := msg.NAck(); nakErr != nil {
				s.logger.Warn().Err(err).Msg("failed to nack a message")
			}
		case <-wCtx.Done():
			return
		}

	}
}
