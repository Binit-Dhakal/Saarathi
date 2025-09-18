package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
	ridematchingv1 "github.com/Binit-Dhakal/Saarathi/trips/tripspb/proto/ride_matching"
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

func RegisterIntegrationEventHandlers(subscriber am.MessageSubsciber, handlers am.MessageHandler) (err error) {
	_, err = subscriber.Subscribe("TRIPS", handlers, am.MessageFilter{
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
	payload := event.Payload().(*ridematchingv1.DriverAccepted)

	acceptedDTO := dto.DriverAccepted{
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
	}
	return h.integrationSvc.DriverAccepted(ctx, acceptedDTO)
}
