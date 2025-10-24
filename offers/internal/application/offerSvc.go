package application

import (
	"context"
	"errors"
	"time"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type OfferService interface {
	CreateTripReadModel(ctx context.Context, payload domain.TripReadModelDTO) error
	ProcessCandidatesList(ctx context.Context, candidates domain.MatchedDriversDTO) error
	HandleNoCandidateFound(ctx context.Context, candidates domain.MatchedDriversDTO) error
	ProcessAcceptedOffer(ctx context.Context, replyDto domain.OfferReplyDTO) error
	ProcessRejectedOffer(ctx context.Context, replyDto domain.OfferReplyDTO) error
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

	err = o.candidateRepo.SaveFirstAttemptUnix(ctx, payload.TripID)
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

func (o *offerSvc) hasExpired(ctx context.Context, tripID string) (bool, error) {
	firstAttemptUnix, err := o.candidateRepo.GetFirstAttemptUnix(ctx, tripID)
	if err != nil {
		return false, err
	}

	return time.Since(time.Unix(firstAttemptUnix, 0)) > 5*time.Minute, nil
}

func (o *offerSvc) ProcessCandidatesList(ctx context.Context, candidates domain.MatchedDriversDTO) error {
	tripID := candidates.TripID

	if err := o.candidateRepo.SaveCandidates(ctx, tripID, candidates.CandidateDriversID); err != nil {
		return err
	}
	return o.tryOfferToCandidates(ctx, candidates, 0)
}

func (o *offerSvc) tryOfferToCandidates(ctx context.Context, candidates domain.MatchedDriversDTO, index int) error {
	tripID := candidates.TripID

	tripDetail, err := o.tripRepo.GetTripDetails(ctx, tripID)
	if err != nil {
		return err
	}

	for {
		driverID, err := o.candidateRepo.GetNextCandidates(ctx, tripID, index)
		if err != nil {
			if errors.Is(err, domain.ErrCandidateListExhausted) {
				firstAttemptUnix, err := o.candidateRepo.GetFirstAttemptUnix(ctx, tripID)
				if err != nil {
					return err
				}

				evt := ddd.NewEvent(domain.RetryMatchingEvent, &domain.RetryMatching{
					SagaID:           tripDetail.SagaID,
					TripID:           tripDetail.TripID,
					CarType:          tripDetail.CarType,
					PickUp:           tripDetail.PickUp,
					DropOff:          tripDetail.DropOff,
					NextAttempt:      candidates.Attempt + 1,
					FirstAttemptUnix: firstAttemptUnix,
				})
				return o.publisher.Publish(ctx, evt)
			}
			return err
		}

		presenceServer, err := o.driverGateway.CheckPresence(ctx, driverID)
		if err != nil {
			index++
			continue
		}

		locked, err := o.driverGateway.TryAcquireLock(ctx, driverID, tripID)
		if err != nil || locked {
			index++
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
}

func (o *offerSvc) ProcessAcceptedOffer(ctx context.Context, replyDto domain.OfferReplyDTO) error {
	evtPayload := &domain.TripOfferAccepted{
		SagaID:   replyDto.SagaID,
		TripID:   replyDto.TripID,
		DriverID: replyDto.DriverID,
	}
	acceptEvt := ddd.NewEvent(domain.TripOfferAcceptedEvent, evtPayload)

	return o.publisher.Publish(ctx, acceptEvt)
}

func (o *offerSvc) HandleNoCandidateFound(ctx context.Context, candidates domain.MatchedDriversDTO) error {
	tripDetail, err := o.tripRepo.GetTripDetails(ctx, candidates.TripID)
	if err != nil {
		return err
	}

	firstAttemptUnix, err := o.candidateRepo.GetFirstAttemptUnix(ctx, candidates.TripID)
	if err != nil {
		return err
	}

	evt := ddd.NewEvent(domain.NoCandidateMatchedEvent, &domain.NoCandidateMatched{
		SagaID:           tripDetail.SagaID,
		TripID:           tripDetail.TripID,
		CarType:          tripDetail.CarType,
		PickUp:           tripDetail.PickUp,
		DropOff:          tripDetail.DropOff,
		NextAttempt:      candidates.Attempt + 1,
		FirstAttemptUnix: firstAttemptUnix,
	})

	return o.publisher.Publish(ctx, evt)
}

// handle case for both rejected and timedout offer
func (o *offerSvc) ProcessRejectedOffer(ctx context.Context, replyDto domain.OfferReplyDTO) error {
	currentIdx, err := o.candidateRepo.IncrementCandidateCounter(ctx, replyDto.TripID)
	if err != nil {
		return err
	}

	err = o.driverGateway.ReleaseLock(ctx, replyDto.DriverID, replyDto.TripID)
	if err != nil {
		return err
	}

	err = o.candidateRepo.AddRejectedDriver(ctx, replyDto.TripID, replyDto.DriverID)
	if err != nil {
		return err
	}

	candidates := domain.MatchedDriversDTO{
		SagaID:  replyDto.SagaID,
		TripID:  replyDto.TripID,
		Attempt: int32(currentIdx),
	}

	return o.tryOfferToCandidates(ctx, candidates, currentIdx)
}
