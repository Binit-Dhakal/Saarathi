package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	CandidatesListKey = "trip:%s:candidates"
	AttemptCounterKey = "trip:%s:counter"
	RedisTTL          = 10 * time.Minute
)

type tripCandidatesRepo struct {
	client *redis.Client
}

func NewTripCandidatesRepo(client *redis.Client) domain.TripCandidatesRepository {
	return &tripCandidatesRepo{
		client: client,
	}
}

func (t *tripCandidatesRepo) SaveCandidates(ctx context.Context, tripID string, driverIDs []string) error {
	pipe := t.client.Pipeline()

	drivers := make([]any, len(driverIDs))
	for i, driverID := range driverIDs {
		drivers[i] = driverID
	}

	listKey := fmt.Sprintf(CandidatesListKey, tripID)
	pipe.Del(ctx, listKey)
	pipe.RPush(ctx, listKey, drivers...)
	pipe.Expire(ctx, listKey, RedisTTL)

	counterKey := fmt.Sprintf(AttemptCounterKey, tripID)
	pipe.Set(ctx, counterKey, 0, RedisTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save candidates in redis pipeline: %w", err)
	}

	return nil
}

func (t *tripCandidatesRepo) IncrementCandidateCounter(ctx context.Context, tripID string) (int, error) {
	key := fmt.Sprintf(AttemptCounterKey, tripID)

	result, err := t.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter: %w", err)
	}

	return int(result), nil
}

func (t *tripCandidatesRepo) GetNextCandidates(ctx context.Context, tripID string, index int) (string, error) {
	key := fmt.Sprintf(CandidatesListKey, tripID)

	driverID, err := t.client.LIndex(ctx, key, int64(index)).Result()
	if err != nil {
		switch err {
		case redis.Nil:
			return "", fmt.Errorf("candidate list exhausted for trip %s", tripID)
		default:
			return "", fmt.Errorf("failed to get candidates at index %d: %w", index, err)
		}
	}

	return driverID, nil
}
