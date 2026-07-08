package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)


type URLStorage interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
}

type urlStorage struct {
	redis *redis.Client
}

func NewURLStorage(redis *redis.Client) URLStorage {
	return &urlStorage{redis: redis}
}

func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	if err := s.redis.Set(context.Background(), code, url, exp*time.Second).Err(); err != nil {
		return err
	}
	return nil
}


var ErrCodeNotFound = errors.New("code doesn't exist")
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	url, err := s.redis.Get(context.Background(), code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}
		return "", err
	}
	return url, nil 
}
