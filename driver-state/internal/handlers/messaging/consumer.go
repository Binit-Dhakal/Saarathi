package messaging

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
	"github.com/Binit-Dhakal/Saarathi/pkg/events"
)

type TripOfferHandler struct {
	notifier domain.DriverNotifier
}

func NewTripOfferHandler(notifier domain.DriverNotifier) *TripOfferHandler {
	return &TripOfferHandler{
		notifier: notifier,
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

	err := t.notifier.NotifyClient(event.DriverID, dto.EventSend{
		Event: "TRIP_OFFER_REQUEST",
		Data:  offerReq,
	})

	if err != nil {
		fmt.Printf("couldn't send to driver %s: %v\n", event.DriverID, err)
		return err
	}

	return nil
}
