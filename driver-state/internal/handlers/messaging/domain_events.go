package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/driverspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type domainHandlers struct {
	publisher am.EventPublisher
}

func NewDomainHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return domainHandlers{publisher: publisher}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.DriverOfferRespondedEvent,
		domain.DriverOfferTimedOutEvent,
	)
}

func (h domainHandlers) HandleEvent(ctx context.Context, event ddd.Event) (err error) {
	switch event.EventName() {
	case domain.DriverOfferRespondedEvent:
		err = h.onOfferResponded(ctx, event)
	case domain.DriverOfferTimedOutEvent:
		err = h.onOfferTimedOut(ctx, event)
	}

	return err
}

func (h domainHandlers) onOfferResponded(ctx context.Context, event ddd.Event) error {
	var err error
	payload := event.Payload().(*domain.Offer)
	switch payload.Status {
	case domain.OfferAccepted:
		replyPayload := ddd.NewEvent(driverspb.OfferAcceptedEvent, driverspb.OfferAccepted{
			DriverId: payload.DriverID,
			TripId:   payload.TripID,
		})

		err = h.publisher.Publish(ctx, driverspb.OfferAcceptedEvent, replyPayload)

	case domain.OfferRejected:
		replyPayload := ddd.NewEvent(driverspb.OfferRejectedEvent, driverspb.OfferRejected{
			DriverId: payload.DriverID,
			TripId:   payload.TripID,
		})

		err = h.publisher.Publish(ctx, driverspb.OfferRejectedEvent, replyPayload)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h domainHandlers) onOfferTimedOut(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(domain.DriverOfferTimeout)

	switch payload.Status {
	case domain.OfferTimedOut:
		replyPayload := ddd.NewEvent(driverspb.OfferTimedoutEvent, driverspb.OfferTimedout{
			TripId: payload.TripID,
		})

		err := h.publisher.Publish(ctx, driverspb.OfferTimedoutEvent, replyPayload)
		if err != nil {
			return err
		}
	}

	return nil
}
