package redis

import (
	"context"
	"time"

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
	pipe := l.client.TxPipeline()

	geoKey := "geo:drivers:available"
	ttlKey := "geo:driver:" + loc.DriverID + ":ttl"

	pipe.GeoAdd(
		context.Background(),
		geoKey,
		&redis.GeoLocation{
			Longitude: loc.Longitude,
			Latitude:  loc.Latitude,
			Name:      loc.DriverID,
		},
	)

	pipe.Set(context.Background(), ttlKey, 1, 30*time.Second)

	_, err := pipe.Exec(context.Background())
	return err
}

func (l *LocationRepo) RemoveActiveGeoLocation(driverID string) error {
	key := "geo:drivers:available"

	return l.client.ZRem(
		context.Background(),
		key,
		driverID,
	).Err()
}
