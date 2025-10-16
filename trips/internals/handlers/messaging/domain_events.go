package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/common"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
)

type domainHandlers struct {
	tripPublisher  am.EventPublisher
	offerPublisher am.EventPublisher
}

func NewDomainEventHandlers(tripPublisher am.EventPublisher, offerPublisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers{
		tripPublisher:  tripPublisher,
		offerPublisher: offerPublisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.TripCreatedEvent,
	)
}

func (h domainHandlers) HandleEvent(ctx context.Context, event ddd.Event) error {
	switch event.EventName() {
	case domain.TripCreatedEvent:
		h.onTripCreated(ctx, event)
	}

	return nil
}

func (h domainHandlers) onTripCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TripCreated)

	createdEvent := &tripspb.TripRequested{
		SagaId:   payload.SagaID,
		TripId:   payload.TripID,
		Distance: payload.Distance,
		Price:    int32(payload.Price),
		PickUp:   &common.Coordinates{Lng: payload.Pickup[0], Lat: payload.Pickup[1]},
		DropOff:  &common.Coordinates{Lng: payload.DropOff[0], Lat: payload.DropOff[1]},
		CarType:  string(payload.CarType),
	}

	evt := ddd.NewEvent(tripspb.TripRequestedEvent, createdEvent)

	err := h.offerPublisher.Publish(ctx, tripspb.TripRequestedEvent, evt)
	if err != nil {
		return err
	}

	err = h.tripPublisher.Publish(ctx, tripspb.TripRequestedEvent, evt)
	if err != nil {
		return err
	}

	return nil
}
