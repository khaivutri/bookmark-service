package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisDB struct {
	client *redis.Client
}

// NewRedisDB constructs a Redis-backed DBCache using the provided client.
func NewRedisDB(client *redis.Client) DBCache {
	return &redisDB{client: client}
}

// SetCacheData stores a value under a group key and field key with TTL.
func (r *redisDB) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error {
	_, err := r.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, cacheGroupKey, cacheKey, value)
		p.Expire(ctx, cacheGroupKey, exp)
		return nil
	})
	return err
}

// GetCacheData retrieves a cached value from the Redis hash stored under the given group key.
func (r *redisDB) GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error) {
	return r.client.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
}

// DeleteCache removes all cached entries under the given group key.
func (r *redisDB) DeleteCache(ctx context.Context, cacheGroupKey string) error {
	return r.client.Del(ctx, cacheGroupKey).Err()
}
