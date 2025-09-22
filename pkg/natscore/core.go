package natscore

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type Broker struct {
	nc      *nats.Conn
	mu      sync.Mutex
	subs    []*nats.Subscription
	logger  zerolog.Logger
	timeout time.Duration
}

var _ am.RequestTransport = (*Broker)(nil)

func NewCoreBroker(nc *nats.Conn, logger zerolog.Logger) *Broker {
	timeout := 5 * time.Second
	return &Broker{
		nc:      nc,
		logger:  logger,
		timeout: timeout,
	}
}

func (b *Broker) Request(ctx context.Context, subject string, msg am.RawMessage) (am.RawMessage, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	resp, err := b.nc.Request(subject, data, b.timeout)
	if err != nil {
		return nil, err
	}
	var reply am.RawMessage
	if err := json.Unmarshal(resp.Data, &reply); err != nil {
		return nil, err
	}

	return reply, nil
}

func (b *Broker) Reply(subject string, handler func(ctx context.Context, req am.RawMessage) (am.RawMessage, error), options ...am.SubscriberOption) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	subCfg := am.NewSubscriberConfig(options)

	var sub *nats.Subscription
	var err error

	// for competing consumer
	if groupName := subCfg.GroupName(); groupName != "" {
		sub, err = b.nc.QueueSubscribe(subject, groupName, b.handleRequest(subCfg, handler))
	} else {
		sub, err = b.nc.Subscribe(subject, b.handleRequest(subCfg, handler))
	}
	if err != nil {
		return err
	}

	b.subs = append(b.subs, sub)
	return nil
}

func (b *Broker) handleRequest(cfg am.SubscriberConfig, handler func(ctx context.Context, req am.RawMessage) (am.RawMessage, error)) func(*nats.Msg) {
	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	return func(msg *nats.Msg) {
		var req rawMessage

		if err := json.Unmarshal(msg.Data, &req); err != nil {
			b.logger.Warn().Err(err).Msg("failed to unmarshal request")
			return
		}

		reply, err := handler(context.Background(), &req)
		if err != nil {
			b.logger.Error().Err(err).Msg("Handler error")
			return
		}

		if msg.Reply != "" && (reply.Data() != nil) {
			respData, err := json.Marshal(reply)
			if err != nil {
				b.logger.Error().Err(err).Msg("failed to marshal reply")
				return
			}
			if err := msg.Respond(respData); err != nil {
				b.logger.Error().Err(err).Msg("failed to send reply")
			}
		}
	}
}

func (b *Broker) Close() (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subs {
		if !sub.IsValid() {
			continue
		}
		err = sub.Drain()
		if err != nil {
			return
		}
	}

	b.subs = nil
	return
}
