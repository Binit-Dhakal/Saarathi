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
	pipe := d.client.Pipeline()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, driverID := range driverIDs {
		pipe.HMGet(ctx, fmt.Sprintf("driver:meta:%s", driverID), "vehicleType")
	}

	responses, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var metadata []domain.DriverVehicleMetadata
	for i, resp := range responses {
		m := domain.DriverVehicleMetadata{
			DriverID:    driverIDs[i],
			VehicleType: resp.String(),
		}

		metadata = append(metadata, m)
	}

	return metadata, nil
}
