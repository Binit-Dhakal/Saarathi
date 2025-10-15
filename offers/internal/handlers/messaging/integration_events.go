package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type integrationHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

func NewIntegrationEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterIntegrationHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) error {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})
	err := subscriber.Subscribe(tripspb.TripRequestedEvent, evtMsgHandler)

	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case tripspb.TripRequestedEvent:
		return h.onTripRequested(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onTripRequested(ctx context.Context, event T) error {
	payload := event.Payload().(*tripspb.TripRequested)

	matchDriversPayload := &offerspb.RideMatchingRequested{
		SagaId:            payload.SagaId,
		TripId:            payload.TripId,
		DropOff:           payload.DropOff,
		PickUp:            payload.PickUp,
		CarType:           payload.CarType,
		MaxSearchRadiusKm: 3,
	}

	matchDriverEvt := ddd.NewEvent(offerspb.RideMatchingRequestedEvent, matchDriversPayload)

	return h.publisher.Publish(ctx, offerspb.RideMatchingRequestedEvent, matchDriverEvt)
}
