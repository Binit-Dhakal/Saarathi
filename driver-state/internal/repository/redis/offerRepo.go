package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/redis/go-redis/v9"
)

type offerRepository struct {
	client *redis.Client
}

func NewOfferRepository(client *redis.Client) domain.OfferRepository {
	return &offerRepository{
		client: client,
	}
}

func offerKey(id string) string {
	return fmt.Sprintf("offer:%s", id)
}

func (r *offerRepository) FindByID(ctx context.Context, id string) (*domain.Offer, error) {
	data, err := r.client.Get(ctx, offerKey(id)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("Redis get failed: %w", err)
	}

	var offer domain.Offer
	if err := json.Unmarshal(data, &offer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal offer: %w", err)
	}

	return &offer, nil
}

func (r *offerRepository) Save(ctx context.Context, offer *domain.Offer) error {
	data, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("failed to marshal offer: %w", err)
	}

	remainingTTL := time.Until(offer.ExpiresAt.Add(5 * time.Minute))
	if remainingTTL <= 0 {
		remainingTTL = 1 * time.Second
	}

	err = r.client.Set(ctx, offerKey(offer.Aggregate.ID()), data, remainingTTL).Err()
	if err != nil {
		return fmt.Errorf("Redis set failed: %w", err)
	}

	return nil
}
