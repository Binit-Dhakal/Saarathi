package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/driverspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type commandHandler struct {
	offerSvc application.OfferService
}

func NewCommandHandler(offerSvc application.OfferService) ddd.CommandHandler {
	return &commandHandler{
		offerSvc: offerSvc,
	}
}

func RegisterCommandHandlers(subscriber am.PublishSubscriber, handlers ddd.CommandHandler) error {
	cmdHandler := func(ctx context.Context, cmd am.IncomingCommandMessage) error {
		_, err := handlers.HandleCommand(ctx, cmd)
		return err
	}

	return subscriber.ReceiveCommand(
		driverspb.CommandChannel,
		cmdHandler,
		am.MessageFilter{driverspb.TripOfferCommand},
	)
}

func (c *commandHandler) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case driverspb.TripOfferCommand:
		return c.doSendTripOffer(ctx, cmd)
	}

	return nil, fmt.Errorf("unsupported command: %s", cmd.CommandName())
}

func (c *commandHandler) doSendTripOffer(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*driverspb.TripOfferRequest)
	fmt.Printf("New trip offer consumed by Driver service: %v\n", payload)

	offer := domain.NewOffer(payload.TripId, payload.DriverId, 15*time.Second, "")
	err := c.offerSvc.SendOffer(&offer)
	if err != nil {
		return ddd.NewReply(driverspb.OfferAckReply, driverspb.OfferAck{
			Accepted: false,
		}), nil

	}

	err = c.offerSvc.SendOffer(&offer)
	if err != nil {
		return ddd.NewReply(driverspb.OfferAckReply, driverspb.OfferAck{
			Accepted: false,
		}), nil
	}

	return ddd.NewReply(driverspb.OfferAckReply, driverspb.OfferAck{
		Accepted: true,
	}), nil
}
