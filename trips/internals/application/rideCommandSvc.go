package application

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
)

type RideCommandService interface {
	AcceptDriverToTrip(ctx context.Context, input dto.AcceptDriver) error
}

type rideCommandService struct {
	tripRepo     domain.TripRepository
	evtPublisher am.EventPublisher
}

var _ RideCommandService = (*rideCommandService)(nil)

func NewRideCommandService(tripRepo domain.TripRepository, bus am.EventPublisher) *rideCommandService {
	return &rideCommandService{
		tripRepo:     tripRepo,
		evtPublisher: bus,
	}
}

func (c *rideCommandService) AcceptDriverToTrip(ctx context.Context, input dto.AcceptDriver) error {
	err := c.tripRepo.AssignDriverToTrip(input.TripID, input.DriverID)
	if err != nil {
		return err
	}

	// TODO: populate driver info:
	// Later after building driver cache repository in trips service(with GRPC fallback)
	// TODO: Outbox transactional pattern
	// remove inconsistency between event publish and database change
	evt := ddd.NewEvent(tripspb.TripConfirmedEvent, &tripspb.TripConfirmed{
		TripId:   input.TripID,
		DriverId: input.DriverID,
	})

	err = c.evtPublisher.Publish(ctx, tripspb.TripAggregateChannel, evt)
	if err != nil {
		return err
	}

	return nil
}
