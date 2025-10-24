package redis

import (
	"context"
	"fmt"

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
		Count:     8,
	}).Val()

	return candidates
}

func (r *rideMatchingRepository) GetRejectedDriver(ctx context.Context, tripID string) ([]string, error) {
	key := fmt.Sprintf("trip:rejected:%s", tripID)
	members, err := r.client.SMembers(ctx, key).Result()
	if err == redis.Nil {
		return []string{}, nil
	}

	if err != nil {
		return nil, err
	}

	return members, nil
}
