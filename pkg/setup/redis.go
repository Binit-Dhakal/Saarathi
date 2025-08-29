package setup

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func SetupRedis(addr string) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:         addr,
		MinIdleConns: 10,
		PoolSize:     20,
		PoolTimeout:  time.Second * 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return client, nil
}
