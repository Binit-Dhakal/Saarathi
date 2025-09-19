package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
)

type integrationHandlers[T ddd.Event] struct {
	integrationSvc application.RideIntegrationService
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(integrationSvc application.RideIntegrationService) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		integrationSvc: integrationSvc,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err = subscriber.Subscribe("TRIPS", evtMsgHandler, am.MessageFilter{
		rmspb.DriverAcceptedEvent,
	})

	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case rmspb.DriverAcceptedEvent:
		return h.onDriverAccepted(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onDriverAccepted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*rmspb.DriverAccepted)

	acceptedDTO := dto.DriverAccepted{
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
	}
	return h.integrationSvc.DriverAccepted(ctx, acceptedDTO)
}
