package application

import (
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
)

type OfferService interface {
	SendTripOffer(driverID string, offer any) error
}

type offerService struct {
	notifier domain.DriverNotifier
}

func NewOfferService(notifier domain.DriverNotifier) OfferService {
	return &offerService{
		notifier: notifier,
	}
}

func (o *offerService) SendTripOffer(driverID string, offer any) error {
	err := o.notifier.NotifyClient(driverID, dto.OfferResponse{
		Event: "NEW_OFFER",
		Data:  offer,
	})

	if err != nil {
		fmt.Printf("couldn't send to driver %s: %v\n", driverID, err)
	}

	return nil
}
