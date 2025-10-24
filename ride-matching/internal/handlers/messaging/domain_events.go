package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
)

type domainHandlers struct {
	publisher am.EventPublisher
}

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.MatchingCandidatesEvent,
		domain.NoDriverAvailableEvent,
	)
}

func (h domainHandlers) HandleEvent(ctx context.Context, event ddd.Event) error {
	switch event.EventName() {
	case domain.MatchingCandidatesEvent:
		return h.onCandidatesMatched(ctx, event)
	case domain.NoDriverAvailableEvent:
		return h.onNoDriverAvailable(ctx, event)
	}

	return nil
}

func (h domainHandlers) onCandidatesMatched(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.MatchingCandidates)

	matchDriversPayload := &rmspb.CandidatesMatched{
		SagaId:           payload.SagaID,
		TripId:           payload.TripID,
		DriverIds:        payload.DriverIds,
		Attempt:          payload.Attempt,
		FirstAttemptUnix: payload.FirstAttemptUnix,
	}

	matchDriverEvt := ddd.NewEvent(rmspb.RMSCandidatesMatchedEvent, matchDriversPayload)

	return h.publisher.Publish(ctx, rmspb.RMSAggregateChannel, matchDriverEvt)
}

func (h domainHandlers) onNoDriverAvailable(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.NoDriverAvailable)

	noDriversPayload := &rmspb.NoDriverMatched{
		TripId:           payload.TripID,
		SagaId:           payload.SagaID,
		Attempt:          payload.Attempt,
		FirstAttemptUnix: payload.FirstAttemptUnix,
	}

	evt := ddd.NewEvent(rmspb.RMSNoDriverMatchedEvent, noDriversPayload)

	return h.publisher.Publish(ctx, rmspb.RMSAggregateChannel, evt)
}
