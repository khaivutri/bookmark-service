package cache

import (
	"context"
	"time"
)

type DBCache interface {

	// SetCacheData sets a cached data item with the given key, value and expiration time.
	// The cache group key is used to categorize the cached data item and can be used to clear all cached data items with the same group key.
	// The cache key is used to uniquely identify the cached data item within its group.
	// The value is the cached data item.
	// The expiration time is the duration after which the cached data item should be automatically deleted.
	SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error

	// GetCacheData retrieves a cached data item with the given group key and key.
	// It returns the cached data item value and an error if the cached data item does not exist.
	// The cache group key is used to categorize the cached data item and can be used to clear all cached data items with the same group key.
	// The cache key is used to uniquely identify the cached data item within its group.
	GetCacheData(ctx context.Context,  cacheGroupKey, cacheKey string) ([]byte, error)

	// DeleteCache deletes a cached data item with the given key.
	// It returns an error if the cached data item does not exist.
	// The key is used to uniquely identify the cached data item.
	DeleteCache(ctx context.Context, cacheGroupKey string) error
}