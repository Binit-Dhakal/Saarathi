package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type OfferService interface {
	ProcessTripOffer(offerID string, driverID string, tripID string, result string) error
	SendOffer(offer *domain.Offer) error
}

type offerService struct {
	publisher ddd.EventPublisher[ddd.Event]
	notifier  domain.DriverNotifier
}

var _ OfferService = (*offerService)(nil)

func NewOfferService(publisher ddd.EventPublisher[ddd.Event], notifier domain.DriverNotifier) OfferService {
	return &offerService{
		publisher: publisher,
		notifier:  notifier,
	}
}

func (o *offerService) ProcessTripOffer(offerID string, driverID string, tripID string, result string) error {
	// TODO: get offer data by searching repo
	offer := domain.NewOffer(tripID, driverID, 15*time.Second, "")

	var event ddd.Event

	switch result {
	case "accepted":
		event, _ = offer.Accept()
	case "rejected":
		event, _ = offer.Reject()
	case "timeout":
		event, _ = offer.TimeOut()
	}

	return o.publisher.Publish(context.Background(), event)
}

func (o *offerService) SendOffer(offer *domain.Offer) error {
	offerReq := dto.OfferRequestDriver{
		TripID:    offer.TripID,
		ExpiresAt: offer.ExpiresAt,
	}
	err := o.notifier.NotifyClient(offer.DriverID, dto.EventSend{
		Event: "TRIP_OFFER_REQUEST",
		Data:  offerReq,
	})

	if err != nil {
		fmt.Printf("couldn't send to driver %s: %v\n", offer.DriverID, err)
		return err
	}

	return nil
}
