package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/redis/go-redis/v9"
)

type presenceRepo struct {
	client *redis.Client
}

func NewPresenceRepo(client *redis.Client) domain.PresenceRepository {
	return &presenceRepo{
		client: client,
	}
}

func (p *presenceRepo) GetDriverInstanceLocation(driverID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	key := fmt.Sprintf("driverID:driver-state-instance:%s", driverID)

	instanceLoc := p.client.Get(ctx, key).Val()
	if instanceLoc == "" {
		return "", fmt.Errorf("Empty instance Location for driverID:%s", driverID)
	}

	return instanceLoc, nil
}
