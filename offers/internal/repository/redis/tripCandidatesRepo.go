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
	FirstAttemptKey   = "trip:%s:first_attempt"
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

func (t *tripCandidatesRepo) SaveFirstAttemptUnix(ctx context.Context, tripID string) error {
	key := fmt.Sprintf(FirstAttemptKey, tripID)
	return t.client.Set(ctx, key, time.Now().Unix(), RedisTTL).Err()
}

func (t *tripCandidatesRepo) GetFirstAttemptUnix(ctx context.Context, tripID string) (int64, error) {
	key := fmt.Sprintf(FirstAttemptKey, tripID)
	ts, err := t.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, fmt.Errorf("no first attempt timestamp found for trip %s", tripID)
		}
		return 0, fmt.Errorf("failed to get first attempt timestamp: %w", err)
	}

	return ts, nil
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
			return "", domain.ErrCandidateListExhausted
		default:
			return "", fmt.Errorf("failed to get candidates at index %d: %w", index, err)
		}
	}

	return driverID, nil
}

func (t *tripCandidatesRepo) AddRejectedDriver(ctx context.Context, tripID string, driverID string) error {
	key := fmt.Sprintf("trip:rejected:%s", tripID)
	if err := t.client.SAdd(ctx, key, driverID).Err(); err != nil {
		return err
	}

	t.client.Expire(ctx, key, RedisTTL)
	return nil
}
