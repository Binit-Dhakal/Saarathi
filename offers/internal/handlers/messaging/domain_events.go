package messaging

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/common"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type domainHandlers struct {
	publisher am.EventPublisher
}

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.RideMatchingInitializedEvent,
		domain.TripOfferEvent,
		domain.TripOfferAcceptedEvent,
	)
}

func (h domainHandlers) HandleEvent(ctx context.Context, event ddd.Event) error {
	switch event.EventName() {
	case domain.RideMatchingInitializedEvent:
		return h.onRideMatchingInitialized(ctx, event)
	case domain.TripOfferEvent:
		return h.onTripOffer(ctx, event)
	case domain.TripOfferAcceptedEvent:
		return h.onTripOfferAccepted(ctx, event)
	}

	return nil
}

func (h domainHandlers) onRideMatchingInitialized(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.RideMatchingInitialized)

	matchDriversPayload := &offerspb.RideMatchingRequested{
		SagaId:            payload.SagaID,
		TripId:            payload.TripID,
		PickUp:            &common.Coordinates{Lng: payload.PickUp[0], Lat: payload.PickUp[1]},
		DropOff:           &common.Coordinates{Lng: payload.DropOff[0], Lat: payload.DropOff[1]},
		CarType:           payload.CarType,
		MaxSearchRadiusKm: 3,
	}

	matchDriverEvt := ddd.NewEvent(offerspb.RideMatchingRequestedEvent, matchDriversPayload)

	fmt.Println("Offer asked the rms to provide driver data", matchDriverEvt)
	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel, matchDriverEvt)
}

func (h domainHandlers) onTripOffer(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TripOffer)

	routingKey := fmt.Sprintf(offerspb.OfferInstanceEventChannel, payload.PresenceServerID)

	p := &offerspb.TripOfferRequested{
		SagaId:   payload.SagaID,
		TripId:   payload.TripID,
		Price:    payload.Price,
		Distance: payload.Distance,
	}

	evt := ddd.NewEvent(offerspb.TripOfferRequestedEvent, p)
	return h.publisher.Publish(ctx, routingKey, evt)
}

func (h domainHandlers) onTripOfferAccepted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TripOfferAccepted)

	p := &offerspb.TripOfferAccepted{
		SagaId:   payload.SagaID,
		DriverId: payload.DriverID,
		TripId:   payload.TripID,
	}

	evt := ddd.NewEvent(offerspb.TripOfferAcceptedEvent, p)
	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel, evt)
}
