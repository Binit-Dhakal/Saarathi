package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/application"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/domain"
)

type integrationHandlers[T ddd.Event] struct {
	svc application.RiderUpdateService
}

func NewIntegrationEventHandlers(svc application.RiderUpdateService) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		svc: svc,
	}
}

func RegisterIntegrationHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) error {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err := subscriber.Subscribe(tripspb.TripAggregateChannel, evtMsgHandler, am.MessageFilter{
		tripspb.TripAssignedEvent,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	var err error
	switch event.EventName() {
	case tripspb.TripAssignedEvent:
		err = h.onTripAssigned(ctx, event)
	}

	return err
}

func (h integrationHandlers[T]) onTripAssigned(ctx context.Context, event T) error {
	payload := event.Payload().(*tripspb.TripAssigned)

	confirmedDTO := &domain.TripConfirmedDTO{
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
		RiderID:  payload.RiderId,
	}

	return h.svc.SendTripCompleteDetail(ctx, confirmedDTO)
}
