package application

import (
	"context"
	"time"

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

	startMatchPayload := &domain.RideMatchingInitialized{
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

	return o.tryOfferToCandidates(ctx, candidates)
}

func (o *offerSvc) tryOfferToCandidates(ctx context.Context, candidates domain.MatchedDriversDTO) error {
	const MaxAttemptWindow = 5 * time.Minute
	tripID := candidates.TripID

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

		if time.Since(time.Unix(candidates.FirstAttemptUnix, 0)) > MaxAttemptWindow {
			break
		}

		locked, err := o.driverGateway.TryAcquireLock(ctx, driverID, tripID)
		if err != nil {
			continue
		}
		if !locked {
			continue
		}

		payload := &domain.TripOffer{
			SagaID:           tripDetail.SagaID,
			TripID:           tripID,
			DriverID:         driverID,
			Price:            tripDetail.Price,
			Distance:         tripDetail.Distance,
			PickUp:           tripDetail.PickUp,
			DropOff:          tripDetail.DropOff,
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

	evt := ddd.NewEvent(domain.NoCandidateMatchedEvent, &domain.NoCandidateMatched{
		SagaID:            tripDetail.SagaID,
		TripID:            tripDetail.TripID,
		CarType:           tripDetail.CarType,
		PickUp:            tripDetail.PickUp,
		DropOff:           tripDetail.DropOff,
		MaxSearchRadiusKm: candidates.SearchRadius + 1,
		Attempt:           candidates.Attempt + 1,
		FirstAttemptUnix:  candidates.FirstAttemptUnix,
	})

	return o.publisher.Publish(ctx, evt)

}

func (o *offerSvc) ProcessAcceptedOffer(ctx context.Context, replyPayload domain.OfferAcceptedReplyDTO) error {
	evtPayload := &domain.TripOfferAccepted{
		SagaID:   replyPayload.SagaID,
		TripID:   replyPayload.TripID,
		DriverID: replyPayload.DriverID,
	}
	acceptEvt := ddd.NewEvent(domain.TripOfferAcceptedEvent, evtPayload)

	return o.publisher.Publish(ctx, acceptEvt)
}
