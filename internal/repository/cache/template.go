package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TemplateCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type templateCache struct {
	redis *redis.Client
}

func NewTemplateCache(redis *redis.Client) TemplateCache {
	return &templateCache{
		redis: redis,
	}
}

func (c *templateCache) Get(ctx context.Context, key string) (string, error) {
	return c.redis.Get(ctx, key).Result()
}

func (c *templateCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.redis.Set(ctx, key, value, expiration).Err()
}

func (c *templateCache) Del(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}
