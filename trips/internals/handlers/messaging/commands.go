package messaging

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
)

type commandHandler struct {
	svc application.RideCommandService
}

func NewCommandHandler(svc application.RideCommandService) ddd.CommandHandler {
	return &commandHandler{
		svc: svc,
	}
}

func RegisterCommandHandlers(subscriber am.CommandSubscriber, handlers ddd.CommandHandler) error {
	cmdHandler := am.CommandMessageHandlerFunc(func(ctx context.Context, cmd am.IncomingCommandMessage) (ddd.Reply, error) {
		return handlers.HandleCommand(ctx, cmd)
	})

	return subscriber.Subscribe(
		tripspb.CommandChannel,
		cmdHandler,
		am.MessageFilter{tripspb.AcceptDriverCommand},
		am.GroupName("trips-commands"),
	)
}

func (c *commandHandler) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case tripspb.AcceptDriverCommand:
		return c.doAcceptDriver(ctx, cmd)
	}

	return nil, fmt.Errorf("unsupported command: %s", cmd.CommandName())
}

func (c *commandHandler) doAcceptDriver(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*tripspb.AcceptDriver)
	acceptedDTO := dto.AcceptDriver{
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
	}

	err := c.svc.AcceptDriverToTrip(ctx, acceptedDTO)
	if err != nil {
		return nil, err
	}

	return ddd.NewReply(cmd.CommandName(), tripspb.AcceptDriverResponse{
		Success: true,
	}), nil
}
