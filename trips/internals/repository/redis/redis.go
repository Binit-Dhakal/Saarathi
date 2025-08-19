package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisFareRepository struct {
	client *redis.Client
}

func NewRedisFareRepository(client *redis.Client) *RedisFareRepository {
	return &RedisFareRepository{
		client: client,
	}
}

func (rf *RedisFareRepository) CreateEphemeralFareEntry(fare *domain.FareQuote) (string, error) {
	id := uuid.New().String()

	data, err := json.Marshal(fare)
	if err != nil {
		return "", err
	}

	key := "fare:" + id
	err = rf.client.Set(context.Background(), key, data, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}

	return id, nil
}

func (rf *RedisFareRepository) GetEphemeralFareEntry(id string) (*domain.FareQuote, error) {
	key := "fare:" + id

	data, err := rf.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var fare domain.FareQuote
	if err := json.Unmarshal([]byte(data), &fare); err != nil {
		return nil, err
	}

	return &fare, nil
}

func (rf *RedisFareRepository) DeleteEphemeralFareEntry(id string) error {
	key := "fare:" + id
	return rf.client.Del(context.Background(), key).Err()
}
