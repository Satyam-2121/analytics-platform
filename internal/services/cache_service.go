package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	Redis *redis.Client
}

func NewCacheService(redisClient *redis.Client) *CacheService {
	return &CacheService{
		Redis: redisClient,
	}
}

func (c *CacheService) Set(
	ctx context.Context,
	key string,
	value string,
) error {

	return c.Redis.Set(
		ctx,
		key,
		value,
		5*time.Minute,
	).Err()
}

func (c *CacheService) Get(
	ctx context.Context,
	key string,
) (string, error) {

	return c.Redis.Get(
		ctx,
		key,
	).Result()
}
