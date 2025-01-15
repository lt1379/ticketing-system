package cache

import (
	"context"
	"github.com/lts1379/ticketing-system/infrastructure/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type ITestCache interface {
	Set(ctx context.Context, key string, value interface{})
	Get(ctx context.Context, key string) (interface{}, error)
}

type TestCache struct {
	RedisClient *redis.Client
}

func NewTestCache(redisClient *redis.Client) ITestCache {
	return &TestCache{RedisClient: redisClient}
}

func (c *TestCache) Set(ctx context.Context, key string, value interface{}) {
	err := c.RedisClient.Set(ctx, key, value, time.Second*30)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while save redis")
	}
}

func (c *TestCache) Get(ctx context.Context, key string) (interface{}, error) {
	return c.RedisClient.Get(ctx, key).Result()
}
