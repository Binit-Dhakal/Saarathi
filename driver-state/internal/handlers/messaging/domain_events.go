package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/driverspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.ReplyPublisher
}

func NewDomainHandlers(publisher am.ReplyPublisher) ddd.EventHandler[ddd.Event] {
	return domainHandlers[ddd.Event]{publisher: publisher}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.DriverOfferRespondedEvent,
		domain.DriverOfferTimedOutEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	switch event.EventName() {
	case domain.DriverOfferRespondedEvent:
		err = h.onOfferResponded(ctx, event)
	case domain.DriverOfferTimedOutEvent:
		err = h.onOfferTimedOut(ctx, event)
	}

	return err
}

func (h domainHandlers[T]) onOfferResponded(ctx context.Context, event T) error {
	var err error
	payload := event.Payload().(*domain.DriverOfferResponded)
	switch payload.Status {
	case domain.OfferAccepted:
		replyPayload := ddd.NewReply(driverspb.OfferAcceptedReply, driverspb.OfferAccepted{
			DriverId: payload.DriverID,
			TripId:   payload.TripID,
		})

		err = h.publisher.Publish(ctx, driverspb.ReplyChannel, replyPayload)

	case domain.OfferRejected:
		replyPayload := ddd.NewReply(driverspb.OfferRejectedReply, driverspb.OfferRejected{
			DriverId: payload.DriverID,
			TripId:   payload.TripID,
		})

		err = h.publisher.Publish(ctx, driverspb.ReplyChannel, replyPayload)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h domainHandlers[T]) onOfferTimedOut(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.DriverOfferTimeout)

	switch payload.Status {
	case domain.OfferTimedOut:
		replyPayload := ddd.NewReply(driverspb.OfferTimedoutReply, driverspb.OfferTimedout{
			TripId: payload.TripID,
		})

		err := h.publisher.Publish(ctx, driverspb.ReplyChannel, replyPayload)
		if err != nil {
			return err
		}
	}

	return nil
}
