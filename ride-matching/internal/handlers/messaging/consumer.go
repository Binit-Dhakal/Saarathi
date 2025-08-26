package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

func (h *TripEventHandler) HandleTripEvent(body []byte) error {
	var event events.TripEventCreated
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("Failed to unmarshal event: %v", err)
	}

	driverCandidates := h.matchingSvc.FindDrivers(event.PickUp[0], event.PickUp[1])
	onlineCandidates := h.driverInfoSvc.GetOnlineDrivers(driverCandidates)

	// fetch metadata
	metadatas, err := h.driverInfoSvc.GetDriversMetadata(onlineCandidates)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// shortlist candidates based on carType
	var shortlistDrivers []domain.DriverVehicleMetadata
	selectedVehicle := strings.ToUpper(event.CarType)
	for _, metadata := range metadatas {
		vt := strings.ToUpper(metadata.VehicleType)
		if vt == selectedVehicle {
			shortlistDrivers = append(shortlistDrivers, metadata)
		}
	}

	for _, driver := range shortlistDrivers {
		instanceID, _ := h.presenceSvc.GetDriverInstance(driver.DriverID)

		routingKey := messagebus.DriverRoutingKey(events.TripEventCreated{}.EventName(), instanceID)

		h.publisher.Publish(context.Background(), messagebus.TripEventsExchange, routingKey, event)
	}

	return nil
}
