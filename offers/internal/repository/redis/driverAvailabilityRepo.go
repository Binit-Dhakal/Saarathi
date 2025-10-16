package redis

import (
	"context"
	"fmt"
	"time"

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

func (d *driverAvailabilityRepo) TryAcquireLock(ctx context.Context, driverID string, tripID string) (bool, error) {
	lockResult, err := d.client.SetNX(ctx, fmt.Sprintf("lock:driver:%s", driverID), tripID, 30*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("redis lock acquisition failed for %s:%w", driverID, err)

	}

	if !lockResult {
		return false, nil // lock was already held
	}

	return true, nil
}

func (d *driverAvailabilityRepo) ReleaseLock(ctx context.Context, driverID string, tripID string) error {
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then 
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)
	lockKey := fmt.Sprintf("lock:driver:%s", driverID)

	result, err := script.Run(ctx, d.client, []string{lockKey}, tripID).Result()
	if err != nil {
		return fmt.Errorf("failed to execute lock release script for %s: %w", driverID, err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("lock not released: mismatch or not held by %s", tripID)
	}

	return nil
}
