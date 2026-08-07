package repository

import (
	"context"
	"errors"
	"time"

	"github.com/khaivutri/bookmark-service/internal/repository/bookmark"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
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
	db    bookmark.BookmarkRepository
}

// NewURLStorage constructs a new URLStorage instance using Redis and PostgreSQL repositories.
func NewURLStorage(redis *redis.Client, db bookmark.BookmarkRepository) URLStorage {
	return &urlStorage{redis: redis, db: db}
}

// StoreURL persists a mapping of short code to original URL in Redis.
func (s *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	return s.redis.Set(ctx, code, url, exp).Err()
}


var ErrCodeNotFound = errors.New("code doesn't exist")


// GetURL retrieves the original URL for a given short code. It first checks Redis
// for short codes of length 7. If not found, it checks PostgreSQL for short
// codes of length 10. If the code is not found in either location, it returns
// ErrCodeNotFound.
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	switch len(code) {
	case 7:
		return s.getFromRedis(ctx, code)
	case 10:
		return s.getFromPostgres(ctx, code)
	default:
		return "", ErrCodeNotFound
	}
}

func (s *urlStorage) getFromRedis(ctx context.Context, code string) (string, error) {
	value, err := s.redis.Get(ctx, code).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}
	return value, err
}

func (s *urlStorage) getFromPostgres(ctx context.Context, code string) (string, error) {
	record, err := s.db.GetBookmarkByCode(ctx, code)
	if errors.Is(err, dbutils.ErrRecordNotFound) {
		return "", ErrCodeNotFound
	} 
	if err != nil {
		return "", err
	}
	return record.URL, nil
}
