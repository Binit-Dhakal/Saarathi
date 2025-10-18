package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
)

type integrationHandlers[T ddd.Event] struct {
	rideSvc application.RideService
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(rideSvc application.RideService) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		rideSvc: rideSvc,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err = subscriber.Subscribe(offerspb.OfferAggregateChannel, evtMsgHandler, am.MessageFilter{
		offerspb.TripOfferAcceptedEvent,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case offerspb.TripOfferAcceptedEvent:
		return h.onDriverAccepted(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onDriverAccepted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*offerspb.TripOfferAccepted)
	acceptedDTO := dto.AcceptDriver{
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
		SagaID:   payload.SagaId,
	}

	return h.rideSvc.AcceptDriverToTrip(ctx, acceptedDTO)
}
