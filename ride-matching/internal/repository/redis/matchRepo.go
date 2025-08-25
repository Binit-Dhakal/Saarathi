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

func NewRideMatchingRepository(client *redis.Client) domain.RedisRideMatchingRepository {
	return &rideMatchingRepository{
		client: client,
	}
}

func (r *rideMatchingRepository) deleteUnavailableDriver(ctx context.Context, expiredDrivers []string) {
	if len(expiredDrivers) > 0 {
		r.client.ZRem(ctx, "geo:drivers:available", expiredDrivers)
	}
}

func (r *rideMatchingRepository) checkDriverAvailabilty(ctx context.Context, candidates []string) ([]string, []string) {
	validDrivers := make([]string, 0, len(candidates))
	expiredDrivers := make([]string, 0)

	pipe := r.client.Pipeline()
	ttlCmds := make([]*redis.IntCmd, len(candidates))

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		fmt.Println(err)
		return []string{}, []string{}
	}

	for i, cmd := range ttlCmds {
		if cmd.Val() > 0 {
			validDrivers = append(validDrivers, candidates[i])
		} else {
			expiredDrivers = append(expiredDrivers, candidates[i])
		}
	}

	return validDrivers, expiredDrivers
}

func (r *rideMatchingRepository) FindNearestDriver(lat, lon float64) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	candidates := r.client.GeoSearch(ctx, "geo:drivers:available", &redis.GeoSearchQuery{
		Latitude:  lat,
		Longitude: lon,
		Radius:    3,
		Sort:      "ASC",
		Count:     50,
	}).Val()

	if len(candidates) == 0 {
		return []string{}
	}

	validDrivers, expiredDrivers := r.checkDriverAvailabilty(ctx, candidates)

	r.deleteUnavailableDriver(ctx, expiredDrivers)

	return validDrivers

}

func (r *rideMatchingRepository) IsDriverAvailable(driverID string) bool {
	value := r.client.Get(context.Background(), fmt.Sprintf("driver:state:%s", driverID))

	if value.Val() == "available" {
		return true
	}

	return false
}
