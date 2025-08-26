package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/events"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
)

type MatchingService interface {
	FindDrivers(lon float64, lat float64) []string
}

type matchingService struct {
	publisher messagebus.Publisher
	matchRepo domain.RedisRideMatchingRepository
	consumer  messagebus.Consumer
}

func NewMatchingService(publisher messagebus.Publisher, matchRepo domain.RedisRideMatchingRepository) MatchingService {
	return &matchingService{
		publisher: publisher,
		matchRepo: matchRepo,
	}
}

// currently our algorithm just find drivers based on geographical location
func (m *matchingService) FindDrivers(lon float64, lat float64) []string {
	return m.matchRepo.FindNearestDriver(lon, lat)
}

func (m *matchingService) HandleNewTripEvent(ctx context.Context, event *events.TripEventCreated) error {
	// find nearest driver - geosearch
	nearDrivers := m.matchRepo.FindNearestDriver(event.PickUp[1], event.PickUp[0])

	metadatas, err := m.getDriverMetadata(nearDrivers)
	if err != nil {
		return err
	}

	var candidateDrivers []domain.DriverVehicleMetadata
	selectedVehicle := strings.ToUpper(event.CarType)
	for _, metadata := range metadatas {
		vt := strings.ToUpper(metadata.VehicleType)
		if vt == selectedVehicle {
			candidateDrivers = append(candidateDrivers, metadata)
		}
	}

	// loop over driver State and send the event if available
	for _, driver := range candidateDrivers {
		// Search in redis if driver:state:{driverId}
		// Search for queue where driverId is connected to: "driver:presence:{driverID}"
		// Topic Exchange to that queue
		// if !m.matchRepo.IsDriverAvailable(driver.DriverID) {
		// 	continue
		// }
		//
		offer := events.TripOffer{
			TripID:    event.RideID,
			DriverID:  driver.DriverID,
			PickUp:    event.PickUp,
			DropOff:   event.DropOff,
			CarType:   selectedVehicle,
			ExpiresAt: time.Now().Add(15 * time.Second),
		}

		fmt.Printf("Published message: %+v\n", offer)
		err = m.publisher.Publish(context.Background(), "trip_offer_exchange", "trip.offer", offer)
		if err != nil {
			continue
		}

		// some goroutine wait for response
	}

	return nil
}

// func (m *matchingService) waitForDriverResponse(ctx context.Context, tripID string, driverID string, timeout time.Duration) (*events.TripOfferResponse, error) {
// 	m.consumer.Consume(context.Background(), "")
//
// }
