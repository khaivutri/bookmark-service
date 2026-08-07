package bookmark_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/khaivutri/bookmark-service/internal/model"
	cacheMocks "github.com/khaivutri/bookmark-service/internal/repository/cache/mock_cache"
	bookmark "github.com/khaivutri/bookmark-service/internal/service/bookmark"
	serviceMocks "github.com/khaivutri/bookmark-service/internal/service/bookmark/mocks"
	"github.com/stretchr/testify/require"
)

const (
	testUserID = "user-1"
	testPage   = 2
	testLimit  = 10
	groupKey   = "get_bookmarks_user-1"
	cacheKey   = "2_10"
)

var (
	errCache   = errors.New("cache error")
	errClear   = errors.New("clear cache error")
	errService = errors.New("service error")
	wantResult = &bookmark.GetBookmarkResponse{
		Bookmarks: []*model.Bookmark{{
			Description: "description", 
			URL: "https://example.com", 
			UserID: testUserID,
		}},
		Total:     1,
	}
)

func TestBookmarkServiceWithCache_GetBookmarks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		cached     []byte
		cacheErr   error
		clearErr   error
		setErr     error
		serviceRes *bookmark.GetBookmarkResponse
		serviceErr error
		want       *bookmark.GetBookmarkResponse
		wantErr    error
	}{
		{
			name: "cache hit", 
			cached: mustJSON(t, wantResult), 
			want: wantResult,
		},
		{
			name: "cache miss", 
			cacheErr: errCache, 
			serviceRes: wantResult, 
			want: wantResult,
		},
		{
			name: "empty cache", 
			serviceRes: wantResult, 
			want: wantResult,
		},
		{
			name: "invalid cache is cleared", 
			cached: []byte("invalid"), 
			serviceRes: wantResult, 
			want: wantResult,
		},
		{
			name: "service error", 
			cacheErr: errCache, 
			serviceErr: errService, 
			wantErr: errService,
		},
		{
			name: "set cache error does not affect result", 
			cacheErr: errCache, 
			setErr: errCache, 
			serviceRes: 
			wantResult, 
			want: wantResult,
		},
		{
			name: "clear cache error does not affect fallback", 
			cached: []byte("invalid"), 
			clearErr: errClear, 
			serviceRes: wantResult, 
			want: wantResult,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			svc := serviceMocks.NewBookmarkService(t)
			cache := cacheMocks.NewDBCache(t)
			setupGetBookmarks(svc, cache, tc)

			got, err := bookmark.NewServiceWithCache(svc, cache).GetBookmarks(ctx, testUserID, testPage, testLimit)

			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestBookmarkServiceWithCache_InvalidatesBeforeMutation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		call     func(bookmark.BookmarkService) error
		setup    func(*serviceMocks.BookmarkService)
		cacheErr error
		wantErr  error
	}{
		{
			name: "add bookmark",
			call: func(s bookmark.BookmarkService) error {
				_, err := s.AddBookmark(context.Background(), "description", "url", testUserID)
				return err
			},
			setup: func(s *serviceMocks.BookmarkService) {
				s.On("AddBookmark", context.Background(), "description", "url", testUserID).Return(nil, nil).Once()
			},
		},
		{
			name: "update bookmark",
			call: func(s bookmark.BookmarkService) error {
				return s.UpdateBookmark(context.Background(), testUserID, "bookmark-1", "description", "url")
			},
			setup: func(s *serviceMocks.BookmarkService) {
				s.On("UpdateBookmark", context.Background(), testUserID, "bookmark-1", "description", "url").Return(nil).Once()
			},
		},
		{
			name: "delete bookmark",
			call: func(s bookmark.BookmarkService) error {
				return s.DeleteBookmark(context.Background(), testUserID, "bookmark-1")
			},
			setup: func(s *serviceMocks.BookmarkService) {
				s.On("DeleteBookmark", context.Background(), testUserID, "bookmark-1").Return(nil).Once()
			},
		},
		{
			name:     "returns clear cache error",
			call:     func(s bookmark.BookmarkService) error { return s.DeleteBookmark(context.Background(), testUserID, "bookmark-1") },
			cacheErr: errClear,
			wantErr:  errClear,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := serviceMocks.NewBookmarkService(t)
			cache := cacheMocks.NewDBCache(t)
			cache.On("DeleteCache", context.Background(), groupKey).Return(tc.cacheErr).Once()
			if tc.setup != nil {
				tc.setup(svc)
			}

			err := tc.call(bookmark.NewServiceWithCache(svc, cache))
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func setupGetBookmarks(svc *serviceMocks.BookmarkService, cache *cacheMocks.DBCache, tc struct {
	name       string
	cached     []byte
	cacheErr   error
	clearErr   error
	setErr     error
	serviceRes *bookmark.GetBookmarkResponse
	serviceErr error
	want       *bookmark.GetBookmarkResponse
	wantErr    error
}) {
	ctx := context.Background()
	cache.On("GetCacheData", ctx, groupKey, cacheKey).Return(tc.cached, tc.cacheErr).Once()

	invalid := tc.cacheErr == nil && len(tc.cached) > 0 && !json.Valid(tc.cached)
	if invalid {
		cache.On("DeleteCache", ctx, groupKey).Return(tc.clearErr).Once()
	}

	if tc.cacheErr != nil || len(tc.cached) == 0 || invalid {
		svc.On("GetBookmarks", ctx, testUserID, testPage, testLimit).Return(tc.serviceRes, tc.serviceErr).Once()
		if tc.serviceErr == nil {
			cache.On("SetCacheData", ctx, groupKey, cacheKey, mustJSONValue(tc.serviceRes), 24*time.Hour).Return(tc.setErr).Once()
		}
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	return mustJSONValue(value)
}

func mustJSONValue(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}