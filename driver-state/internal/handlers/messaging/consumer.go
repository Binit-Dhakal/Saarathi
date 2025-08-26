package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
	"github.com/Binit-Dhakal/Saarathi/pkg/events"
)

type TripOfferHandler struct {
	notifyClient domain.DriverNotifier
}

func NewTripOfferHandler(notifyClient domain.DriverNotifier) *TripOfferHandler {
	return &TripOfferHandler{
		notifyClient: notifyClient,
	}
}

func (t *TripOfferHandler) HandleOfferRequest(ctx context.Context, evt events.Event) error {
	event := evt.(*events.TripOfferRequest)

	offerReq := dto.OfferRequestDriver{
		TripID:    event.TripID,
		PickUp:    event.PickUp,
		DropOff:   event.DropOff,
		ExpiresAt: event.ExpiresAt,
	}

	req := dto.EventSend{
		Event: "TRIP_OFFER_REQUEST",
		Data:  offerReq,
	}

	err := t.notifyClient.NotifyClient(event.DriverID, req)
	if err != nil {
		return err
	}

	return nil
}
