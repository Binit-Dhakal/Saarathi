package command

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
)

type Broker interface {
	Request(ctx context.Context, subject string, msg am.RawMessage) (am.RawMessage, error)
	Responder
}
