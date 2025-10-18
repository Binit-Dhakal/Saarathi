package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/redis/go-redis/v9"
)

const projectionKeyPrefix = "projection:trip:%s"

type tripProjectionRepository struct {
	client *redis.Client
}

var _ domain.TripProjectionRepository = (*tripProjectionRepository)(nil)

func NewTripProjectionRepository(client *redis.Client) domain.TripProjectionRepository {
	return &tripProjectionRepository{
		client: client,
	}
}

func (t *tripProjectionRepository) SetTripPayload(ctx context.Context, tripID string, payload []byte, expiration time.Duration) error {
	key := fmt.Sprintf(projectionKeyPrefix, tripID)

	err := t.client.Set(ctx, key, payload, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save trip projection payload to redis for trip %s: %w", tripID, err)
	}

	return nil
}
