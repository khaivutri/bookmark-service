package bookmark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository/cache"
	"github.com/rs/zerolog/log"
)

const (
	getBookmarksCacheGroupKeyFormat = "get_bookmarks_%s"
	getBookmarksCacheKeyFormat      = "%d_%d"
	getBookmarksCacheExpTime        = 24 * time.Hour
)

type bookmarkServiceWithCache struct {
	svc BookmarkService
	c   cache.DBCache
}

// NewServiceWithCache returns a BookmarkService decorated with cache invalidation and retrieval.
func NewServiceWithCache(s BookmarkService, c cache.DBCache) BookmarkService {
	return &bookmarkServiceWithCache{
		svc: s,
		c:   c,
	}
}

// AddBookmark invalidates user bookmark cache before adding a new bookmark.
func (s *bookmarkServiceWithCache) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	// create cache key
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)

	// delete cache
	err := s.c.DeleteCache(ctx, cacheGroupKey)
	if err != nil {
		return nil, err
	}

	return s.svc.AddBookmark(ctx, description, url, userID)

}
// GetBookmarks returns cached bookmarks when available, otherwise loads from the underlying service and caches the result.
func (s *bookmarkServiceWithCache) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarkResponse, error) {
	// create cache key
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)

	cacheKey := fmt.Sprintf(getBookmarksCacheKeyFormat, page, limit)

	// get cache data
	cacheData, err := s.c.GetCacheData(ctx, cacheGroupKey, cacheKey)
	if err == nil && len(cacheData) > 0 {
		result := &GetBookmarkResponse{}

		err := json.Unmarshal(cacheData, result)

		if err != nil {

			// error data - > delete cache
			err := s.c.DeleteCache(ctx, cacheGroupKey)

			// log error
			if err != nil {
				log.Err(err).Str("key", cacheGroupKey).Msg("fail to delete cache")
			}

		} else {
			return result, nil
		}
	}

	// if cache not existed, call service
	result, err := s.svc.GetBookmarks(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	// save cache
	resultBytes, err := json.Marshal(result)
	if err == nil {
		cacheErr := s.c.SetCacheData(ctx, cacheGroupKey, cacheKey, resultBytes, getBookmarksCacheExpTime)

		// log error
		if cacheErr != nil {
			log.Err(cacheErr).Str("key", cacheGroupKey).Msg("fail to set cache")
		}
	}

	// return result
	return result, nil
}

// UpdateBookmark invalidates user bookmark cache before updating a bookmark.
func (s *bookmarkServiceWithCache) UpdateBookmark(ctx context.Context, userID, bookmarkID, description, url string) error {
	// create cache key
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)

	// delete cache
	err := s.c.DeleteCache(ctx, cacheGroupKey)
	if err != nil {
		return err
	}

	return s.svc.UpdateBookmark(ctx, userID, bookmarkID, description, url)

}

// DeleteBookmark invalidates user bookmark cache before deleting a bookmark.
func (s *bookmarkServiceWithCache) DeleteBookmark(ctx context.Context, userID, bookmarkID string) error {
	// create cache key
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)

	// delete cache
	err := s.c.DeleteCache(ctx, cacheGroupKey)
	if err != nil {
		return err
	}

	return s.svc.DeleteBookmark(ctx, userID, bookmarkID)
}
