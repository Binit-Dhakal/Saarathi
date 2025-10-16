package messaging

import (
	"context"

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
		domain.RideMatchingInitializedEvent)
}

func (h domainHandlers) HandleEvent(ctx context.Context, event ddd.Event) error {
	switch event.EventName() {
	case domain.RideMatchingInitializedEvent:
		h.onRideMatchingInitialized(ctx, event)
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

	return h.publisher.Publish(ctx, offerspb.RideMatchingRequestedEvent, matchDriverEvt)
}
