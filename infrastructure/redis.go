package infrastructure

import (
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/redis/go-redis/v9"
)
func createRedis() *redis.Client{
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	return redisClient
}