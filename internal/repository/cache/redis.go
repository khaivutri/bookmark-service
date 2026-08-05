package cache

import (
	"github.com/redis/go-redis/v9"
	"context"
	"time"
)

type redisDB struct {
	client *redis.Client
}

func NewRedisDB(client *redis.Client) DBCache{
	return &redisDB{client: client}
}


func (r *redisDB) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error {
	_, err := r.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		r.client.HSet(ctx, cacheGroupKey, cacheKey, value )
		r.client.Expire(ctx, cacheGroupKey, exp)
		return nil
	})
	return err
}
 
func (r *redisDB) GetCacheData(ctx context.Context,  cacheGroupKey, cacheKey string) ([]byte, error) {
	return r.client.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
}

func (r *redisDB) DeleteCache(ctx context.Context, cacheGroupKey string) error {
	return r.client.HDel(ctx, cacheGroupKey).Err()
}
