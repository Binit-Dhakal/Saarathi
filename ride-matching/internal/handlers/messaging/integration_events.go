package messaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/offerspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/rmspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/dto"
)

type integrationHandlers[T ddd.Event] struct {
	matchingSvc   application.MatchingService
	driverInfoSvc application.DriverInfoService
	presenceSvc   application.PresenceService
	publisher     am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(matchingSvc application.MatchingService, driverInfoSvc application.DriverInfoService, presenceSvc application.PresenceService, publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		matchingSvc:   matchingSvc,
		driverInfoSvc: driverInfoSvc,
		presenceSvc:   presenceSvc,
		publisher:     publisher,
	}
}

func RegisterIntegrationEventHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) (err error) {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err = subscriber.Subscribe(offerspb.RideMatchingRequestedEvent, evtMsgHandler, am.GroupName("RMS-Service"))

	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case offerspb.RideMatchingRequestedEvent:
		return h.onMatchingRequest(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onMatchingRequest(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*offerspb.RideMatchingRequested)
	evt := dto.TripCreated{
		SagaID:       payload.GetSagaId(),
		TripID:       payload.GetTripId(),
		PickUp:       payload.GetPickUp(),
		DropOff:      payload.GetDropOff(),
		CarType:      payload.GetCarType(),
		SearchRadius: payload.GetMaxSearchRadiusKm(),
	}

	driverCandidates := h.matchingSvc.FindDrivers(evt.PickUp.GetLng(), evt.PickUp.GetLat())
	onlineCandidates := h.driverInfoSvc.GetOnlineDrivers(driverCandidates)

	// fetch metadata
	metadatas, err := h.driverInfoSvc.GetDriversMetadata(onlineCandidates)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// shortlist candidates based on carType
	var shortlistDrivers []string
	selectedVehicle := strings.ToUpper(payload.CarType)
	for _, metadata := range metadatas {
		vt := strings.ToUpper(metadata.VehicleType)
		if vt == selectedVehicle {
			shortlistDrivers = append(shortlistDrivers, metadata.DriverID)
		}
	}

	replyPayload := &rmspb.MatchingCandidates{
		SagaId:    payload.SagaId,
		TripId:    payload.TripId,
		DriverIds: shortlistDrivers,
	}
	matchedEvt := ddd.NewEvent(rmspb.RMSMatchingCandidatesEvent, replyPayload)
	err = h.publisher.Publish(ctx, rmspb.RMSMatchingCandidatesEvent, matchedEvt)
	if err != nil {
		return err
	}
	return nil
}
