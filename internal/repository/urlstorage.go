package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// URLStorage defines the interface for storing and retrieving URLs.
type URLStorage interface {
	// StoreURL persists a mapping of short code to original URL in Redis.
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	// GetURL retrieves the original URL for a given short code.
	GetURL(ctx context.Context, code string) (string, error)
}


type urlStorage struct {
	redis *redis.Client
}

// NewURLStorage constructs a new URLStorage instance using a Redis client.
func NewURLStorage(redis *redis.Client) URLStorage {
	return &urlStorage{redis: redis}
}

// StoreURL persists a mapping of short code to original URL in Redis.
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	if err := s.redis.Set(context.Background(), code, url, exp*time.Second).Err(); err != nil {
		return err
	}
	return nil
}


// GetURL retrieves the original URL for a given short code.
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
