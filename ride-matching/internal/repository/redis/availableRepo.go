package redis

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/redis/go-redis/v9"
)

type driverAvailableRepository struct {
	client *redis.Client
}

func NewDriverAvailableRepo(client *redis.Client) domain.DriverAvailabilityRepository {
	return &driverAvailableRepository{
		client: client,
	}
}

// TODO
func (d *driverAvailableRepository) IsDriverFree(driverID string) bool {
	// key := "geo:driver:" + driverID + ":ttl"
	// res := d.client.Exists(context.Background(), key).Val()
	// if res == 1 {
	// 	return true
	// }
	//
	return true
}

func (d *driverAvailableRepository) DeleteUnavailableDrivers(expiredDrivers []string) {
	if len(expiredDrivers) > 0 {
		d.client.ZRem(context.Background(), "geo:drivers:available", expiredDrivers)
	}
}

func (d *driverAvailableRepository) BulkCheckDriversOnline(driversID []string) ([]string, []string) {
	ctx := context.Background()

	validDrivers := make([]string, 0, len(driversID))
	expiredDrivers := make([]string, 0)

	pipe := d.client.Pipeline()
	ttlCmds := make([]*redis.IntCmd, len(driversID))

	for i, driverID := range driversID {
		ttlKey := "geo:driver:" + driverID + ":ttl"
		ttlCmds[i] = pipe.Exists(ctx, ttlKey)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		fmt.Println(err)
		return []string{}, []string{}
	}

	for i, cmd := range ttlCmds {
		if cmd.Val() > 0 {
			validDrivers = append(validDrivers, driversID[i])
		} else {
			expiredDrivers = append(expiredDrivers, driversID[i])
		}
	}

	return validDrivers, expiredDrivers
}
