package redis

import (
	"context"
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
func (r *rideMatchingRepository) FindNearestDriver(lon, lat float64) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	candidates := r.client.GeoSearch(ctx, "geo:drivers:available", &redis.GeoSearchQuery{
		Latitude:  lat,
		Longitude: lon,
		Radius:    3,
		Sort:      "ASC",
		Count:     50,
	}).Val()

	return candidates
}
