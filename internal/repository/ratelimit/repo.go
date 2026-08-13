package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RateLimiter interface {
	IncreaseRateLimit(ctx context.Context, key string, exp time.Duration)
	GetCurrentRateLimit(ctx context.Context, key string ) ( int, error)
}

type redisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) RateLimiter {
	return &redisRepo{client: client}
}

func (r *redisRepo) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) {
	_, err := r.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Incr(ctx, key)
		p.Expire(ctx, key, exp)
		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("fail to increase rate limit")
	}
}

func (r *redisRepo) GetCurrentRateLimit(ctx context.Context, key string) (int, error) {
	return r.client.Get(ctx, key).Int()
}