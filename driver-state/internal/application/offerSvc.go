package application

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type OfferService interface {
	ProcessTripOffer(ctx context.Context, offerID string, result string) error
	CreateAndSendOffer(ctx context.Context, offerDto *dto.OfferRequestedDTO) error
	SetNotifier(notifier domain.DriverNotifier)
}

type offerService struct {
	publisher ddd.EventPublisher[ddd.Event]
	repo      domain.OfferRepository
	notifier  domain.DriverNotifier
}

var _ OfferService = (*offerService)(nil)

func NewOfferService(publisher ddd.EventPublisher[ddd.Event], notifier domain.DriverNotifier, repo domain.OfferRepository) OfferService {
	return &offerService{
		publisher: publisher,
		notifier:  notifier,
		repo:      repo,
	}
}

func (o *offerService) ProcessTripOffer(ctx context.Context, offerID string, result string) error {
	// TODO: get offer data by searching repo
	offer, err := o.repo.FindByID(ctx, offerID)
	if err != nil {
		// TODO: handle bad case if other error rather than not found offer
		fmt.Printf("offerID not found: %v", err)
		return nil
	}

	var event ddd.Event

	switch result {
	case "accepted":
		event, err = offer.Accept()
	case "rejected":
		event, err = offer.Reject()
	case "timeout":
		event, err = offer.TimeOut()
	default:
		return fmt.Errorf("invalid offer processing result: %s", result)
	}

	if err != nil {
		return fmt.Errorf("failed to transition offer state: %w", err)
	}

	err = o.repo.Save(ctx, offer)
	if err != nil {
		return fmt.Errorf("failed to saved updated offer state: %w", err)
	}

	return o.publisher.Publish(context.Background(), event)
}

func (o *offerService) CreateAndSendOffer(ctx context.Context, offerDto *dto.OfferRequestedDTO) error {
	offer := domain.NewOffer(offerDto.TripID, offerDto.SagaID, offerDto.DriverID, offerDto.Price, offerDto.Distance)

	err := o.repo.Save(ctx, &offer)
	if err != nil {
		return fmt.Errorf("Failed to save new error: %w", err)
	}

	offerReq := dto.OfferRequestDriver{
		OfferID:   offer.ID(),
		TripID:    offer.TripID,
		ExpiresAt: offer.ExpiresAt,
	}
	err = o.notifier.NotifyClient(offer.DriverID, dto.EventSend{
		Event: "TRIP_OFFER_REQUEST",
		Data:  offerReq,
	})

	if err != nil {
		fmt.Printf("couldn't send to driver %s: %v\n", offer.DriverID, err)
		return err
	}

	return nil
}

func (o *offerService) SetNotifier(notifier domain.DriverNotifier) {
	o.notifier = notifier
}
