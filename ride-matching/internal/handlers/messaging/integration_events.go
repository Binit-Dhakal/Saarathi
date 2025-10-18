package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/dto"
)

type integrationHandlers[T ddd.Event] struct {
	matchingSvc application.MatchingService
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(matchingSvc application.MatchingService) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		matchingSvc: matchingSvc,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err = subscriber.Subscribe(offerspb.OfferAggregateChannel, evtMsgHandler, am.GroupName("RMS-Service"), am.MessageFilter{
		offerspb.RideMatchingRequestedEvent,
	})

	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case offerspb.RideMatchingRequestedEvent:
		return h.onMatchingRequest(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onMatchingRequest(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*offerspb.RideMatchingRequested)
	evt := dto.TripCreated{
		SagaID:       payload.GetSagaId(),
		TripID:       payload.GetTripId(),
		PickUp:       payload.GetPickUp(),
		DropOff:      payload.GetDropOff(),
		CarType:      payload.GetCarType(),
		SearchRadius: payload.GetMaxSearchRadiusKm(),
	}

	return h.matchingSvc.ProcessMatchingRequest(ctx, evt)
}
