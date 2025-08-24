package application

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
)

type MatchingService interface {
	HandleNewTripEvent(ctx context.Context, event *events.TripEventCreated) error
}

type matchingService struct {
	publisher messagebus.Publisher
	matchRepo domain.RideMatchingRepository
	consumer  messagebus.Consumer
}

func NewMatchingService(publisher messagebus.Publisher, matchRepo domain.RideMatchingRepository) MatchingService {
	return &matchingService{
		publisher: publisher,
		matchRepo: matchRepo,
	}
}

func (m *matchingService) HandleNewTripEvent(ctx context.Context, event *events.TripEventCreated) error {
	// find nearest driver - geosearch
	nearDrivers := m.matchRepo.FindNearestDriver(event.PickUp[1], event.PickUp[0])

	//bulk search for metadata of nearDrivers
	metadatas, err := m.matchRepo.BulkSearchDriverMeta(nearDrivers)
	if err != nil {
		return err
	}

	var candidateDrivers []domain.DriverVehicleMetadata
	selectedVehicle := event.CarType
	for _, metadata := range metadatas {
		if metadata.VehicleType == selectedVehicle {
			candidateDrivers = append(candidateDrivers, metadata)
		}
	}

	// loop over driver State and send the event if available
	for _, driver := range candidateDrivers {
		// Search in redis if driver:state:{driverId}
		// Search for queue where driverId is connected to: "driver:presence:{driverID}"
		// Direct Exchange to that queue
		if !m.matchRepo.IsDriverAvailable(driver.DriverID) {
			continue
		}

		offer := events.TripOffer{
			TripID:    event.RideID,
			DriverID:  driver.DriverID,
			PickUp:    event.PickUp,
			DropOff:   event.DropOff,
			CarType:   selectedVehicle,
			ExpiresAt: time.Now().Add(15 * time.Second),
		}

		err = m.publisher.Publish(context.Background(), "trip_offer_exchange", "trip.offer", offer)
		if err != nil {
			continue
		}
	}

	return nil
}

// func (m *matchingService) waitForDriverResponse(ctx context.Context, tripID string, driverID string, timeout time.Duration) (*events.TripOfferResponse, error) {
// 	m.consumer.Consume(context.Background(), "")
//
// }
