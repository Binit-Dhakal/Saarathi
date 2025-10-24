package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/common"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/offerspb"
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
		domain.NoCandidateMatchedEvent,
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
	case domain.NoCandidateMatchedEvent:
		return h.onNoCandidateMatched(ctx, event)
	}

	return nil
}

func (h domainHandlers) onRideMatchingInitialized(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.RideMatchingInitialized)

	matchDriversPayload := &offerspb.RideMatchingRequested{
		SagaId:           payload.SagaID,
		TripId:           payload.TripID,
		PickUp:           &common.Coordinates{Lng: payload.PickUp[0], Lat: payload.PickUp[1]},
		DropOff:          &common.Coordinates{Lng: payload.DropOff[0], Lat: payload.DropOff[1]},
		CarType:          payload.CarType,
		Attempt:          1,
		FirstAttemptUnix: time.Now().Unix(),
	}

	matchDriverEvt := ddd.NewEvent(offerspb.RideMatchingRequestedEvent, matchDriversPayload)

	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel, matchDriverEvt)
}

func (h domainHandlers) onTripOffer(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.TripOffer)

	routingKey := fmt.Sprintf(offerspb.OfferInstanceEventChannel, payload.PresenceServerID)

	p := &offerspb.TripOfferRequested{
		SagaId:   payload.SagaID,
		TripId:   payload.TripID,
		Price:    payload.Price,
		PickUp:   &common.Coordinates{Lng: payload.PickUp[0], Lat: payload.PickUp[1]},
		DropOff:  &common.Coordinates{Lng: payload.DropOff[0], Lat: payload.DropOff[1]},
		DriverId: payload.DriverID,
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

func (h domainHandlers) onNoCandidateMatched(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.NoCandidateMatched)
	delay := min(time.Duration(5*(payload.NextAttempt))*time.Second, 20*time.Second)

	firstAttempt := time.Unix(payload.FirstAttemptUnix, 0)
	expiry := firstAttempt.Add(5 * time.Minute)
	nextAttemptTime := time.Now().Add(delay)

	if nextAttemptTime.After(expiry) {
		fmt.Printf("Trip %s expired after 5 minutes; stopping matching.\n", payload.TripID)

		driverNotFoundPayload := &offerspb.NoDriverFound{
			SagaId:        payload.SagaID,
			TripId:        payload.TripID,
			ExpiredAtUnix: time.Now().Unix(),
		}

		driverNotFoundEvt := ddd.NewEvent(offerspb.NoDriverFoundEvent, driverNotFoundPayload)

		return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel, driverNotFoundEvt)
	}

	time.Sleep(delay)

	fmt.Printf("Retrying match for trip %s after %v\n", payload.TripID, delay)
	matchDriversPayload := &offerspb.RideMatchingRequested{
		SagaId:           payload.SagaID,
		TripId:           payload.TripID,
		PickUp:           &common.Coordinates{Lng: payload.PickUp[0], Lat: payload.PickUp[1]},
		DropOff:          &common.Coordinates{Lng: payload.DropOff[0], Lat: payload.DropOff[1]},
		CarType:          payload.CarType,
		Attempt:          payload.NextAttempt,
		FirstAttemptUnix: payload.FirstAttemptUnix,
	}
	matchDriverEvt := ddd.NewEvent(offerspb.RideMatchingRequestedEvent, matchDriversPayload)

	return h.publisher.Publish(ctx, offerspb.OfferAggregateChannel, matchDriverEvt)
}
