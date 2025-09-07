package setup

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func SetupRedis(addr string) (*goredis.Client, error) {
	opt, err := goredis.ParseURL(addr)
	if err != nil {
		return nil, err
	}
	opt.MinIdleConns = 10
	opt.PoolSize = 20
	opt.PoolTimeout = time.Second * 5

	client := goredis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return client, nil
}
