package repository

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrDependencyDown = errors.New("one or more dependencies are down")
type RedisPinger struct {
	client *redis.Client
}

func NewRedisPinger(client *redis.Client) *RedisPinger {
	return &RedisPinger{client: client}
}

func (p *RedisPinger) Ping(ctx context.Context) error {
	return p.client.Ping(ctx).Err()
}