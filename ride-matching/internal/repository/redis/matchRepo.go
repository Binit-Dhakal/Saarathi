package redis

import (
	"context"

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
func (r *rideMatchingRepository) FindNearestDriver(ctx context.Context, lon, lat, radius float64) []string {
	candidates := r.client.GeoSearch(ctx, "geo:drivers:available", &redis.GeoSearchQuery{
		Latitude:  lat,
		Longitude: lon,
		Radius:    radius,
		Sort:      "ASC",
		Count:     5,
	}).Val()

	return candidates
}
