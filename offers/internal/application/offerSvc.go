package application

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type OfferService interface {
	CreateTripReadModel(ctx context.Context, payload domain.TripReadModelDTO) error
	ProcessCandidatesList(ctx context.Context, candidates domain.MatchedDriversDTO) error
	ProcessAcceptedOffer(ctx context.Context, event domain.OfferAcceptedReplyDTO) error
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

func (o *offerSvc) ProcessCandidatesList(ctx context.Context, candidates domain.MatchedDriversDTO) error {
	tripID := candidates.TripID

	if err := o.candidateRepo.SaveCandidates(ctx, tripID, candidates.CandidateDriversID); err != nil {
		return err
	}

	tripDetail, err := o.tripRepo.GetTripDetails(ctx, tripID)
	if err != nil {
		return err
	}

	for i := range candidates.CandidateDriversID {
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

		payload := domain.TripOffer{
			SagaID:           tripDetail.SagaID,
			TripID:           tripID,
			Price:            tripDetail.Price,
			Distance:         tripDetail.Distance,
			PresenceServerID: presenceServer,
		}

		evt := ddd.NewEvent(domain.TripOfferEvent, payload)

		err = o.publisher.Publish(ctx, evt)
		if err != nil {
			o.driverGateway.ReleaseLock(ctx, driverID, tripID)
			continue
		}

		return nil
	}

	return fmt.Errorf("Temp: Not matched")
}

func (o *offerSvc) ProcessAcceptedOffer(ctx context.Context, replyPayload domain.OfferAcceptedReplyDTO) error {
	evtPayload := domain.TripOfferAccepted{
		SagaID:   replyPayload.SagaID,
		TripID:   replyPayload.TripID,
		DriverID: replyPayload.DriverID,
	}
	acceptEvt := ddd.NewEvent(domain.TripOfferAcceptedEvent, evtPayload)

	return o.publisher.Publish(ctx, acceptEvt)
}
