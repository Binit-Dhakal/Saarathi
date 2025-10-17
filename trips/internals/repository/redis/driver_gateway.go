package redis

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/redis/go-redis/v9"
)

type presenceGatewayRepo struct {
	client *redis.Client
}

var _ domain.PresenceGatewayRepository = (*presenceGatewayRepo)(nil)

func NewPresenceGatewayRepository(client *redis.Client) domain.PresenceGatewayRepository {
	return &presenceGatewayRepo{
		client: client,
	}
}

func (g *presenceGatewayRepo) GetDriverLocation(ctx context.Context, driverID string) (*domain.Coordinate, error) {
	positions, err := g.client.GeoPos(ctx, "geo:drivers:available", driverID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve driver location from redis GeoSet: %w", err)
	}

	if len(positions) == 0 || positions[0] == nil {
		return nil, fmt.Errorf("driver %s location not found in GeoSet  (driver likely offline)", driverID)
	}

	pos := positions[0]

	return &domain.Coordinate{
		Lat: pos.Latitude,
		Lon: pos.Longitude,
	}, nil
}
