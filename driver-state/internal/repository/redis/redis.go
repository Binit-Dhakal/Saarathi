package redis

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/redis/go-redis/v9"
)

type LocationRepo struct {
	client *redis.Client
}

func NewLocationRepo(client *redis.Client) *LocationRepo {
	return &LocationRepo{
		client: client,
	}
}

func (l *LocationRepo) SaveActiveGeoLocation(loc *domain.DriverLocation) error {
	key := "geo:drivers:available"

	return l.client.GeoAdd(
		context.Background(),
		key,
		&redis.GeoLocation{
			Longitude: loc.Longitude,
			Latitude:  loc.Latitude,
			Name:      loc.DriverID,
		},
	).Err()
}

func (l *LocationRepo) RemoveActiveGeoLocation(driverID string) error {
	key := "geo:drivers:available"

	return l.client.ZRem(
		context.Background(),
		key,
		driverID,
	).Err()
}
