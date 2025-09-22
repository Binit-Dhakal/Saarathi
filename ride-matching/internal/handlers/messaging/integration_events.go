package messaging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/driverspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type integrationHandlers[T ddd.Event] struct {
	matchingSvc   application.MatchingService
	driverInfoSvc application.DriverInfoService
	presenceSvc   application.PresenceService
	publisher     am.PublishSender
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(matchingSvc application.MatchingService, driverInfoSvc application.DriverInfoService, presenceSvc application.PresenceService, publisher am.PublishSender) ddd.EventHandler[ddd.Event] {
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

	err = subscriber.Subscribe(tripspb.TripAggregateChannel, evtMsgHandler, am.MessageFilter{
		tripspb.TripCreatedEvent,
	})

	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case tripspb.TripCreatedEvent:
		return h.onTripCreated(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onTripCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*tripspb.TripCreated)
	evt := dto.TripCreated{
		TripID:   payload.GetTripId(),
		Distance: payload.GetDistance(),
		Price:    payload.GetPrice(),
		PickUp:   payload.GetPickUp(),
		DropOff:  payload.GetDropOff(),
		CarType:  payload.GetCarType(),
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
	var shortlistDrivers []domain.DriverVehicleMetadata
	selectedVehicle := strings.ToUpper(payload.CarType)
	for _, metadata := range metadatas {
		vt := strings.ToUpper(metadata.VehicleType)
		if vt == selectedVehicle {
			shortlistDrivers = append(shortlistDrivers, metadata)
		}
	}

	for _, driver := range shortlistDrivers {
		instanceID, _ := h.presenceSvc.GetDriverInstance(driver.DriverID)

		offerEvent := driverspb.TripOfferRequest{
			TripId:    evt.TripID,
			DriverId:  driver.DriverID,
			PickUp:    evt.PickUp,
			DropOff:   evt.DropOff,
			ExpiresAt: timestamppb.New(time.Now().Add(15 * time.Second)),
		}

		cmd := ddd.NewCommand(driverspb.TripOfferCommand, &offerEvent)

		routingKey := fmt.Sprintf(driverspb.CommandChannel, instanceID)

		err = h.publisher.SendCommand(context.Background(), routingKey, cmd)
		if err != nil {
			continue
		}
	}

	return nil
}
