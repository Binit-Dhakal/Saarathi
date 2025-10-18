package messaging

import (
	"context"
	"fmt"
	"os"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
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
	hostName, _ := os.Hostname()

	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err := subscriber.Subscribe(fmt.Sprintf(offerspb.OfferInstanceEventChannel, hostName), evtMsgHandler)
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case offerspb.TripOfferRequestedEvent:
		return h.onOfferRequested(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onOfferRequested(ctx context.Context, event T) error {
	payload := event.Payload().(*offerspb.TripOfferRequested)

	offerRequestedDTO := &dto.OfferRequestedDTO{
		TripID:   payload.TripId,
		SagaID:   payload.SagaId,
		Distance: payload.Distance,
		Price:    payload.Price,
		DriverID: payload.DriverId,
	}

	return h.offerSvc.CreateAndSendOffer(ctx, offerRequestedDTO)
}
