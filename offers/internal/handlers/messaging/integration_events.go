package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/application"
	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/driverspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/tripspb"
)

type integrationHandlers[T ddd.Event] struct {
	offerSvc application.OfferService
}

func NewIntegrationEventHandlers(offerSvc application.OfferService) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		offerSvc: offerSvc,
	}
}

func RegisterIntegrationHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) error {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err := subscriber.Subscribe(tripspb.TripAggregateChannel, evtMsgHandler, am.MessageFilter{
		tripspb.TripRequestedEvent,
	}, am.GroupName("offers-trips-requested"))
	if err != nil {
		return err
	}

	err = subscriber.Subscribe(rmspb.RMSAggregateChannel, evtMsgHandler, am.GroupName("offers-rms-matched"))
	if err != nil {
		return err
	}

	err = subscriber.Subscribe(driverspb.DriverAggregateChannel, evtMsgHandler, am.MessageFilter{
		driverspb.OfferAcceptedEvent,
		driverspb.OfferRejectedEvent,
		driverspb.OfferTimedoutEvent,
	}, am.GroupName("offers-drivers-response"))
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case tripspb.TripRequestedEvent:
		return h.onTripRequested(ctx, event)
	case rmspb.RMSCandidatesMatchedEvent:
		return h.onCandidatesList(ctx, event)
	case rmspb.RMSNoDriverMatchedEvent:
		return h.onNoCandidateFound(ctx, event)
	case driverspb.OfferAcceptedEvent:
		return h.onOfferAccepted(ctx, event)
	case driverspb.OfferRejectedEvent:
		return h.onOfferRejected(ctx, event)
	case driverspb.OfferTimedoutEvent:
		return h.onOfferTimedout(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onTripRequested(ctx context.Context, event T) error {
	payload := event.Payload().(*tripspb.TripRequested)

	readDTO := domain.TripReadModelDTO{
		SagaID:   payload.GetSagaId(),
		TripID:   payload.GetTripId(),
		PickUp:   [2]float64{payload.GetPickUp().GetLng(), payload.GetPickUp().GetLat()},
		DropOff:  [2]float64{payload.GetDropOff().GetLng(), payload.GetDropOff().GetLat()},
		CarType:  payload.GetCarType(),
		Price:    payload.GetPrice(),
		Distance: payload.GetDistance(),
	}

	return h.offerSvc.CreateTripReadModel(ctx, readDTO)
}

func (h integrationHandlers[T]) onCandidatesList(ctx context.Context, event T) error {
	payload := event.Payload().(*rmspb.CandidatesMatched)

	candidatesDTO := domain.MatchedDriversDTO{
		SagaID:             payload.GetSagaId(),
		TripID:             payload.GetTripId(),
		CandidateDriversID: payload.GetDriverIds(),
		Attempt:            payload.GetAttempt(),
	}

	return h.offerSvc.ProcessCandidatesList(ctx, candidatesDTO)
}

func (h integrationHandlers[T]) onOfferAccepted(ctx context.Context, event T) error {
	payload := event.Payload().(*driverspb.OfferAccepted)

	replyDTO := domain.OfferReplyDTO{
		SagaID:   payload.SagaId,
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
		Status:   "accepted",
	}

	return h.offerSvc.ProcessAcceptedOffer(ctx, replyDTO)
}

func (h integrationHandlers[T]) onNoCandidateFound(ctx context.Context, event T) error {
	payload := event.Payload().(*rmspb.NoDriverMatched)

	candidatesDTO := domain.MatchedDriversDTO{
		TripID:             payload.TripId,
		SagaID:             payload.SagaId,
		Attempt:            payload.Attempt,
		CandidateDriversID: make([]string, 0),
	}

	return h.offerSvc.HandleNoCandidateFound(ctx, candidatesDTO)
}

func (h integrationHandlers[T]) onOfferRejected(ctx context.Context, event T) error {
	payload := event.Payload().(*driverspb.OfferRejected)

	replyDTO := domain.OfferReplyDTO{
		SagaID:   payload.SagaId,
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
		Status:   "rejected",
	}

	return h.offerSvc.ProcessAcceptedOffer(ctx, replyDTO)
}

func (h integrationHandlers[T]) onOfferTimedout(ctx context.Context, event T) error {
	payload := event.Payload().(*driverspb.OfferTimedOut)

	replyDTO := domain.OfferReplyDTO{
		SagaID:   payload.SagaId,
		TripID:   payload.TripId,
		DriverID: payload.DriverId,
		Status:   "timedout",
	}

	return h.offerSvc.ProcessAcceptedOffer(ctx, replyDTO)
}
