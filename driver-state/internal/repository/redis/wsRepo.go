package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type wsRepo struct {
	client   *redis.Client
	hostname string
}

func NewWSRepo(client *redis.Client) *wsRepo {
	hostname, _ := os.Hostname()
	return &wsRepo{
		client:   client,
		hostname: hostname,
	}
}

func (w *wsRepo) SaveWSDetail(driverID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	key := fmt.Sprintf("driverID:driver-state-instance:%s", driverID)

	return w.client.Set(ctx, key, w.hostname, 30*time.Second).Err()
}

func (w *wsRepo) DeleteWSDetail(driverID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	key := fmt.Sprintf("driverID:driver-state-instance:%s", driverID)

	return w.client.Del(ctx, key).Err()
}
