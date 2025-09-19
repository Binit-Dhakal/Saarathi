package natscore

import (
	"context"
	"sync"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/command"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type CommandBus struct {
	nc      *nats.Conn
	mu      sync.Mutex
	logger  zerolog.Logger
	timeout time.Duration
}

var _ command.Broker = (*CommandBus)(nil)

func NewCommandBus(nc *nats.Conn, logger zerolog.Logger) *CommandBus {
	timeout := 5 * time.Second
	return &CommandBus{
		nc:      nc,
		logger:  logger,
		timeout: timeout,
	}
}

func (c *CommandBus) Request(ctx context.Context, subject string, msg am.RawMessage) (am.RawMessage, error) {
	resp, err := c.nc.Request(subject, msg.Data(), c.timeout)
	if err != nil {
		return am.NewRawMessage("", "", []byte{}), err
	}
	return am.NewRawMessage(msg.ID(), subject, resp.Data), nil
}

// func (c *CommandBus) SubscribeCommand(subject string, handler func(ctx context.Context, cmd ))
