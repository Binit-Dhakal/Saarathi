package redis

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/redis/go-redis/v9"
)

type driverAvailabilityRepo struct {
	client *redis.Client
}

func NewDriverAvailabilityRepo(client *redis.Client) domain.DriverAvailabilityRepo {
	return &driverAvailabilityRepo{
		client: client,
	}
}

func (d *driverAvailabilityRepo) CheckPresence(ctx context.Context, driverID string) (string, error) {
	key := fmt.Sprintf("driverID:driver-state-instance:%s", driverID)

	instanceLoc := d.client.Get(ctx, key).Val()
	if instanceLoc == "" {
		return "", fmt.Errorf("Empty instance Location for driverID:%s", driverID)
	}

	return instanceLoc, nil
}
