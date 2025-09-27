package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type offerIntentHandler[T ddd.Event] struct {
	offerSvc application.OfferService
}

func NewOfferIntentHandler(offerSvc application.OfferService) ddd.EventHandler[ddd.Event] {
	return offerIntentHandler[ddd.Event]{offerSvc: offerSvc}
}

func RegisterOfferIntentHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.AcceptOfferIntent,
		domain.RejectOfferIntent,
		domain.TimeoutOfferIntent,
	)
}

func (h offerIntentHandler[T]) HandleEvent(ctx context.Context, event T) (err error) {
	switch event.EventName() {
	case domain.AcceptOfferIntent:
		err = h.onOfferAccepted(ctx, event)
	case domain.RejectOfferIntent:
		err = h.onOfferRejected(ctx, event)
	}

	return err
}

func (h offerIntentHandler[T]) onOfferAccepted(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.AcceptOffer)
	return h.offerSvc.ProcessTripOffer(payload.OfferID, payload.DriverID, payload.TripID, "accepted")
}

func (h offerIntentHandler[T]) onOfferRejected(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.RejectOffer)
	return h.offerSvc.ProcessTripOffer(payload.OfferID, payload.DriverID, payload.TripID, "rejected")
}

