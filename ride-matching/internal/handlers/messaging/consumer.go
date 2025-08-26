package messaging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
)

type TripEventHandler struct {
	matchingSvc   application.MatchingService
	driverInfoSvc application.DriverInfoService
	presenceSvc   application.PresenceService
	publisher     messagebus.Publisher
}

func NewTripEventHandler(matchingSvc application.MatchingService, driverInfoSvc application.DriverInfoService, presenceSvc application.PresenceService, publisher messagebus.Publisher) *TripEventHandler {
	return &TripEventHandler{
		matchingSvc:   matchingSvc,
		driverInfoSvc: driverInfoSvc,
		presenceSvc:   presenceSvc,
		publisher:     publisher,
	}
}

func (h *TripEventHandler) HandleTripEvent(ctx context.Context, evt events.Event) error {
	event := evt.(*events.TripEventCreated)
	fmt.Println("Got event: ", event)

	driverCandidates := h.matchingSvc.FindDrivers(event.PickUp[0], event.PickUp[1])
	onlineCandidates := h.driverInfoSvc.GetOnlineDrivers(driverCandidates)
	fmt.Println(driverCandidates, onlineCandidates)

	// fetch metadata
	metadatas, err := h.driverInfoSvc.GetDriversMetadata(onlineCandidates)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(metadatas)

	// shortlist candidates based on carType
	var shortlistDrivers []domain.DriverVehicleMetadata
	selectedVehicle := strings.ToUpper(event.CarType)
	for _, metadata := range metadatas {
		vt := strings.ToUpper(metadata.VehicleType)
		if vt == selectedVehicle {
			shortlistDrivers = append(shortlistDrivers, metadata)
		}
	}

	fmt.Println("Shortlisted: ", shortlistDrivers)

	for _, driver := range shortlistDrivers {
		instanceID, _ := h.presenceSvc.GetDriverInstance(driver.DriverID)
		fmt.Println(instanceID)

		offerEvent := events.TripOfferRequest{
			TripID:    event.RideID,
			DriverID:  driver.DriverID,
			PickUp:    event.PickUp,
			DropOff:   event.DropOff,
			CarType:   selectedVehicle,
			ExpiresAt: time.Now().Add(15 * time.Second),
		}

		routingKey := messagebus.DriverRoutingKey(offerEvent.EventName(), instanceID)

		h.publisher.Publish(context.Background(), messagebus.TripOfferExchange, routingKey, offerEvent)
	}

	return nil
}
