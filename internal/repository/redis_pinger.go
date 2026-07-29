package repository

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// ErrDependencyDown is returned when one or more dependencies are down
var ErrDependencyDown = errors.New("one or more dependencies are down")
// RedisPinger is a test double for the redis client
type RedisPinger struct {
	client *redis.Client
}

// NewRedisPinger constructs a new RedisPinger instance.
func NewRedisPinger(client *redis.Client) *RedisPinger {
	return &RedisPinger{client: client}
}

// Ping checks the connectivity to the Redis client.
func (p *RedisPinger) Ping(ctx context.Context) error {
	return p.client.Ping(ctx).Err()
}