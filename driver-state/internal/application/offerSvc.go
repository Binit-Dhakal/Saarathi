package application

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
)

type OfferService interface {
	SendTripAssignedEvent(driverID string, tripID string, result string) error
}

type offerService struct {
	publisher messagebus.Publisher
}

func NewOfferService(publisher messagebus.Publisher) OfferService {
	return &offerService{
		publisher: publisher,
	}
}

func (o *offerService) SendTripAssignedEvent(driverID string, tripID string, result string) error {
	resp := events.TripOfferResponse{
		TripID:   tripID,
		DriverID: driverID,
		Result:   result,
		AtUnix:   time.Now(),
	}

	err := o.publisher.Publish(
		context.Background(),
		messagebus.TripOfferExchange,
		events.EventOfferResponse,
		resp,
	)

	if err != nil {
		return err
	}

	return nil
}
