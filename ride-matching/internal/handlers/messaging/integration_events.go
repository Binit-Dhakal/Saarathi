package messaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/dto"
)

type integrationHandlers[T ddd.Event] struct {
	matchingSvc   application.MatchingService
	driverInfoSvc application.DriverInfoService
	presenceSvc   application.PresenceService
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(matchingSvc application.MatchingService, driverInfoSvc application.DriverInfoService, presenceSvc application.PresenceService) ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{
		matchingSvc:   matchingSvc,
		driverInfoSvc: driverInfoSvc,
		presenceSvc:   presenceSvc,
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
	createdDTO := dto.TripCreated{
		TripID:   payload.GetTripId(),
		Distance: payload.GetDistance(),
		Price:    payload.GetPrice(),
		PickUp:   payload.GetPickUp(),
		DropOff:  payload.GetDropOff(),
		CarType:  payload.GetCarType(),
	}

	driverCandidates := h.matchingSvc.FindDrivers(createdDTO.PickUp.GetLng(), createdDTO.PickUp.GetLat())
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

	// for _, driver := range shortlistDrivers {
	// 	instanceID, _ := h.presenceSvc.GetDriverInstance(driver.DriverID)
	//
	// 	offerEvent := events.TripOfferRequest{
	// 		TripID:    event.RideID,
	// 		DriverID:  driver.DriverID,
	// 		PickUp:    event.PickUp,
	// 		DropOff:   event.DropOff,
	// 		CarType:   selectedVehicle,
	// 		ExpiresAt: time.Now().Add(15 * time.Second),
	// 	}
	//
	// 	routingKey := messagebus.DriverRoutingKey(events.EventOfferRequest, instanceID)
	//
	// 	fmt.Println("Publishing ", offerEvent, "to", routingKey)
	// 	h.publisher.Publish(context.Background(), messagebus.TripOfferExchange, routingKey, offerEvent)
	// }
	//
	return nil
}
