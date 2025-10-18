package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/rider/internal/domain"
	"github.com/redis/go-redis/v9"
)

type tripPayloadRepository struct {
	client *redis.Client
}

const projectionKeyPrefix = "projection:trip:%s"

var _ domain.TripPayloadRepository = (*tripPayloadRepository)(nil)

func NewTripPayloadRepository(client *redis.Client) domain.TripPayloadRepository {
	return &tripPayloadRepository{
		client: client,
	}
}

func (r *tripPayloadRepository) GetTripFullPayload(ctx context.Context, tripID string) ([]byte, error) {
	key := fmt.Sprintf(projectionKeyPrefix, tripID)

	bytesPayload, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		switch err {
		case redis.Nil:
			return nil, fmt.Errorf("trip projection data not found for trip %s: %w", tripID, err)
		default:
			return nil, fmt.Errorf("failed to retrieve trip projection from redis for trip %s: %w", tripID, err)
		}
	}

	return bytesPayload, nil
}
