package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/env"
	goredis "github.com/redis/go-redis/v9"
)

func SetupRedis() (*goredis.Client, error) {
	redisHost := env.GetEnvWithDefault("REDIS_HOST", "localhost")
	redisPort := env.GetEnvWithDefault("REDIS_PORT", "6379")

	client := goredis.NewClient(&goredis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisHost, redisPort),
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
