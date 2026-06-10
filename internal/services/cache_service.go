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
	key string,
	value string,
) error {

	return c.Redis.Set(
		context.Background(),
		key,
		value,
		5*time.Minute,
	).Err()
}

func (c *CacheService) Get(
	key string,
) (string, error) {

	return c.Redis.Get(
		context.Background(),
		key,
	).Result()
}