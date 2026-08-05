package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/khaivutri/bookmark-service/internal/model"
	bookmarkMocks "github.com/khaivutri/bookmark-service/internal/repository/bookmark/mocks"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const (
	redisCode    = "abc1234"
	postgresCode = "abc1234567"
	testURL      = "https://google.com"
)

func TestURLStorage_StoreURL(t *testing.T) {
	t.Parallel()

	t.Run("stores URL in Redis", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		client := redisPkg.InitMockRedis(t)
		storage := NewURLStorage(client, nil)

		err := storage.StoreURL(ctx, "test", testURL, time.Hour)

		assert.NoError(t, err)
		storedURL, getErr := client.Get(ctx, "test").Result()
		assert.NoError(t, getErr)
		assert.Equal(t, testURL, storedURL)
	})

	t.Run("returns Redis connection error", func(t *testing.T) {
		t.Parallel()
		client := redisPkg.InitMockRedis(t)
		_ = client.Close()
		storage := NewURLStorage(client, nil)

		err := storage.StoreURL(t.Context(), "test", testURL, time.Hour)

		assert.ErrorIs(t, err, redis.ErrClosed)
	})
}

func TestURLStorage_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		code        string
		setup       func(*testing.T, context.Context, *redis.Client, *bookmarkMocks.BookmarkRepository)
		expectedURL string
		expectedErr error
	}{
		{
			name: "gets 7-character code from Redis",
			code: redisCode,
			setup: func(t *testing.T, ctx context.Context, client *redis.Client, _ *bookmarkMocks.BookmarkRepository) {
				assert.NoError(t, client.Set(ctx, redisCode, testURL, time.Hour).Err())
			},
			expectedURL: testURL,
		},
		{
			name:        "maps missing Redis key to ErrCodeNotFound",
			code:        redisCode,
			expectedErr: ErrCodeNotFound,
		},
		{
			name: "returns Redis connection error",
			code: redisCode,
			setup: func(_ *testing.T, _ context.Context, client *redis.Client, _ *bookmarkMocks.BookmarkRepository) {
				_ = client.Close()
			},
			expectedErr: redis.ErrClosed,
		},
		{
			name: "gets 10-character code from PostgreSQL",
			code: postgresCode,
			setup: func(_ *testing.T, ctx context.Context, _ *redis.Client, repo *bookmarkMocks.BookmarkRepository) {
				repo.On("GetBookmarkByCode", ctx, postgresCode).
					Return(&model.Bookmark{URL: testURL}, nil).Once()
			},
			expectedURL: testURL,
		},
		{
			name: "maps missing PostgreSQL record to ErrCodeNotFound",
			code: postgresCode,
			setup: func(_ *testing.T, ctx context.Context, _ *redis.Client, repo *bookmarkMocks.BookmarkRepository) {
				repo.On("GetBookmarkByCode", ctx, postgresCode).
					Return(nil, dbutils.ErrRecordNotFound).Once()
			},
			expectedErr: ErrCodeNotFound,
		},
		{
			name: "returns PostgreSQL error",
			code: postgresCode,
			setup: func(_ *testing.T, ctx context.Context, _ *redis.Client, repo *bookmarkMocks.BookmarkRepository) {
				repo.On("GetBookmarkByCode", ctx, postgresCode).
					Return(nil, errors.New("database unavailable")).Once()
			},
			expectedErr: errors.New("database unavailable"),
		},
		{
			name:        "rejects code with unsupported length",
			code:        "invalid",
			expectedErr: ErrCodeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			client := redisPkg.InitMockRedis(t)
			repo := bookmarkMocks.NewBookmarkRepository(t)
			if tc.setup != nil {
				tc.setup(t, ctx, client, repo)
			}

			link, err := NewURLStorage(client, repo).GetURL(ctx, tc.code)

			assert.Equal(t, tc.expectedURL, link)
			if tc.expectedErr != nil && tc.expectedErr != ErrCodeNotFound {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
