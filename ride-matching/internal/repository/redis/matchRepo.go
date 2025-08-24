package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/redis/go-redis/v9"
)

type rideMatchingRepository struct {
	client *redis.Client
}

func NewRideMatchingRepository(client *redis.Client) domain.RideMatchingRepository {
	return &rideMatchingRepository{
		client: client,
	}
}

func (r *rideMatchingRepository) FindNearestDriver(lat, lon float64) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return r.client.GeoSearch(ctx, "geo:drivers:available", &redis.GeoSearchQuery{
		Latitude:  lat,
		Longitude: lon,
		Radius:    3,
		Sort:      "ASC",
		Count:     50,
	}).Val()
}

func (r *rideMatchingRepository) BulkSearchDriverMeta(driverIDs []string) ([]domain.DriverVehicleMetadata, error) {
	pipe := r.client.Pipeline()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, driverID := range driverIDs {
		pipe.HMGet(ctx, fmt.Sprintf("driver:meta:%s", driverID), "vehicleType")
	}

	responses, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var metadata []domain.DriverVehicleMetadata
	for i, resp := range responses {
		m := domain.DriverVehicleMetadata{
			DriverID:    driverIDs[i],
			VehicleType: resp.String(),
		}

		metadata = append(metadata, m)
	}

	return metadata, nil
}

func (r *rideMatchingRepository) IsDriverAvailable(driverID string) bool {
	value := r.client.Get(context.Background(), fmt.Sprintf("driver:state:%s", driverID))

	if value.Val() == "available" {
		return true
	}

	return false

}
