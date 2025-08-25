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
	HandleNewTripEvent(ctx context.Context, event *events.TripEventCreated) error
}

type matchingService struct {
	publisher     messagebus.Publisher
	matchRepo     domain.RedisRideMatchingRepository
	redisMetaRepo domain.RedisMetaRepository
	pgMetaRepo    domain.PGMetaRepository
	consumer      messagebus.Consumer
}

func NewMatchingService(publisher messagebus.Publisher, matchRepo domain.RedisRideMatchingRepository, redisMetaRepo domain.RedisMetaRepository, pgMetaRepo domain.PGMetaRepository) MatchingService {
	return &matchingService{
		publisher:     publisher,
		matchRepo:     matchRepo,
		redisMetaRepo: redisMetaRepo,
		pgMetaRepo:    pgMetaRepo,
	}
}

func (m *matchingService) getDriverMetadata(driversIDs []string) ([]domain.DriverVehicleMetadata, error) {
	//bulk search for metadata of nearDrivers
	metadatas, err := m.redisMetaRepo.BulkSearchDriverMeta(driversIDs)
	if err != nil {
		return nil, err
	}

	var missing []string
	for _, d := range metadatas {
		if d.VehicleType == "" {
			missing = append(missing, d.DriverID)
		}
	}

	if len(missing) > 0 {
		go func(missing []string) {
			dbMeta, err := m.pgMetaRepo.BulkSearchMeta(driversIDs)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = m.redisMetaRepo.BulkInsertDriverMeta(dbMeta)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(missing)
	}

	return metadatas, nil
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
