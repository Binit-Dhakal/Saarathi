package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/redis/go-redis/v9"
)

type driverMetaRepository struct {
	client *redis.Client
}

func NewCacheDriverMetaRepo(client *redis.Client) domain.RedisMetaRepository {
	return &driverMetaRepository{
		client: client,
	}
}

func (d *driverMetaRepository) BulkInsertDriverMeta(metas []domain.DriverVehicleMetadata) error {
	pipe := d.client.Pipeline()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, meta := range metas {
		key := fmt.Sprintf("driver:meta:%s", meta.DriverID)
		pipe.HSet(ctx, key, "vehicleType", meta.VehicleType)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *driverMetaRepository) BulkSearchDriverMeta(driverIDs []string) ([]domain.DriverVehicleMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmds, err := d.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, driverID := range driverIDs {
			p.HGet(ctx, fmt.Sprintf("driver:meta:%s", driverID), "vehicleType")
		}
		return nil
	})

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var metadata []domain.DriverVehicleMetadata
	for i, cmd := range cmds {
		s, err := cmd.(*redis.StringCmd).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		m := domain.DriverVehicleMetadata{
			DriverID:    driverIDs[i],
			VehicleType: s,
		}

		metadata = append(metadata, m)
	}

	return metadata, nil
}
