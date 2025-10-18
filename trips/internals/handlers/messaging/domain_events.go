package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/common"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/application"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
)

type domainHandlers struct {
	publisher     am.EventPublisher
	projectionSvc application.ProjectionService
}

func NewDomainEventHandlers(publisher am.EventPublisher, projectionSvc application.ProjectionService) ddd.EventHandler[ddd.Event] {
	return &domainHandlers{
		publisher:     publisher,
		projectionSvc: projectionSvc,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.TripCreatedEvent,
		domain.TripMatchedEvent,
	)
}

func (h domainHandlers) HandleEvent(ctx context.Context, event ddd.Event) error {
	switch event.EventName() {
	case domain.TripCreatedEvent:
		return h.onTripCreated(ctx, event)
	case domain.TripMatchedEvent:
		return h.onTripMatched(ctx, event)
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

	err := h.publisher.Publish(ctx, tripspb.TripAggregateChannel, evt)

	if err != nil {
		return err
	}

	return nil
}

func (h domainHandlers) onTripMatched(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TripMatched)

	err := h.projectionSvc.ProjectTripDetails(ctx, payload.TripID, payload.DriverID, payload.RiderID)
	if err != nil {
		return err
	}

	p := &tripspb.TripAssigned{
		SagaId:   payload.SagaID,
		TripId:   payload.TripID,
		DriverId: payload.DriverID,
		RiderId:  payload.RiderID,
	}
	assignedEvt := ddd.NewEvent(tripspb.TripAssignedEvent, p)
	err = h.publisher.Publish(ctx, tripspb.TripAggregateChannel, assignedEvt)
	if err != nil {
		return err
	}

	return nil
}
