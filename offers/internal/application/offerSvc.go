package application

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type OfferService interface {
	CreateTripReadModel(ctx context.Context, payload domain.TripReadModelDTO) error
	ProcessCandidatesList(ctx context.Context, candidates *rmspb.MatchingCandidates) error
}

type offerSvc struct {
	candidateRepo domain.TripCandidatesRepository
	driverGateway domain.DriverAvailabilityRepo
	tripRepo      domain.TripReadModelRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewService(repo domain.TripCandidatesRepository, driverGate domain.DriverAvailabilityRepo, tripRepo domain.TripReadModelRepository, publisher ddd.EventPublisher[ddd.Event]) OfferService {
	return &offerSvc{
		candidateRepo: repo,
		driverGateway: driverGate,
		tripRepo:      tripRepo,
		publisher:     publisher,
	}
}

func (o *offerSvc) CreateTripReadModel(ctx context.Context, payload domain.TripReadModelDTO) error {
	// TODO: need to handle case for tripRepo to have duplicate
	err := o.tripRepo.SaveTrip(ctx, payload)
	if err != nil {
		return err
	}

	startMatchPayload := domain.RideMatchingInitialized{
		SagaID:  payload.SagaID,
		TripID:  payload.TripID,
		CarType: payload.CarType,
		PickUp:  payload.PickUp,
		DropOff: payload.DropOff,
	}
	evt := ddd.NewEvent(domain.RideMatchingInitializedEvent, startMatchPayload)

	return o.publisher.Publish(ctx, evt)
}

func (o *offerSvc) ProcessCandidatesList(ctx context.Context, candidates *rmspb.MatchingCandidates) error {
	tripID := candidates.GetTripId()

	if err := o.candidateRepo.SaveCandidates(ctx, tripID, candidates.GetDriverIds()); err != nil {
		return err
	}

	for i := 0; i < len(candidates.GetDriverIds()); i++ {
		driverID, err := o.candidateRepo.GetNextCandidates(ctx, tripID, i)
		if err != nil {
			continue
		}

		presenceServer, err := o.driverGateway.CheckPresence(ctx, driverID)
		if err != nil {
			continue
		}

		locked, err := o.driverGateway.TryAcquireLock(ctx, driverID, tripID)
		if err != nil {
			continue
		}

		if !locked {
			continue
		}

	}
}
