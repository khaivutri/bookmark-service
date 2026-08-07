package cache_test

import (
	"testing"
	"time"

	cacheRepo "github.com/khaivutri/bookmark-service/internal/repository/cache"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testGroupKey = "bookmarks"
	testCacheKey = "user-1"
)

func newTestCache(t *testing.T) (cacheRepo.DBCache, *redis.Client) {
	t.Helper()

	client := redisPkg.InitMockRedis(t)
	t.Cleanup(func() { _ = client.Close() })

	return cacheRepo.NewRedisDB(client), client
}

func TestNewRedisDB(t *testing.T) {
	t.Parallel()

	db, _ := newTestCache(t)

	assert.NotNil(t, db)
}

func TestRedisDB_SetCacheData(t *testing.T) {
	t.Parallel()

	db, client := newTestCache(t)
	ctx := t.Context()
	value := []byte(`{"id":1}`)
	expiration := time.Minute

	require.NoError(t, db.SetCacheData(ctx, testGroupKey, testCacheKey, value, expiration))

	stored, err := client.HGet(ctx, testGroupKey, testCacheKey).Bytes()
	require.NoError(t, err)
	assert.Equal(t, value, stored)

	ttl, err := client.TTL(ctx, testGroupKey).Result()
	require.NoError(t, err)
	assert.Greater(t, ttl, time.Duration(0))
	assert.LessOrEqual(t, ttl, expiration)
}

func TestRedisDB_SetCacheData_ReturnsErrorWhenClientIsClosed(t *testing.T) {
	t.Parallel()

	db, client := newTestCache(t)
	require.NoError(t, client.Close())

	err := db.SetCacheData(t.Context(), testGroupKey, testCacheKey, []byte("value"), time.Minute)

	assert.Error(t, err)
}

func TestRedisDB_GetCacheData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(t *testing.T, client *redis.Client)
		want    []byte
		wantErr bool
	}{
		{
			name: "returns stored value",
			prepare: func(t *testing.T, client *redis.Client) {
				require.NoError(t, client.HSet(t.Context(), testGroupKey, testCacheKey, []byte("value")).Err())
			},
			want: []byte("value"),
		},
		{
			name:    "returns redis nil when field does not exist",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, client := newTestCache(t)
			if tc.prepare != nil {
				tc.prepare(t, client)
			}

			got, err := db.GetCacheData(t.Context(), testGroupKey, testCacheKey)

			if tc.wantErr {
				assert.ErrorIs(t, err, redis.Nil)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRedisDB_DeleteCache(t *testing.T) {
	t.Parallel()

	db, client := newTestCache(t)
	ctx := t.Context()
	require.NoError(t, client.HSet(ctx, testGroupKey, testCacheKey, []byte("value")).Err())

	require.NoError(t, db.DeleteCache(ctx, testGroupKey))

	_, err := client.HGet(ctx, testGroupKey, testCacheKey).Result()
	assert.ErrorIs(t, err, redis.Nil)
}

func TestRedisDB_DeleteCache_ReturnsErrorWhenClientIsClosed(t *testing.T) {
	t.Parallel()

	db, client := newTestCache(t)
	require.NoError(t, client.Close())

	err := db.DeleteCache(t.Context(), testGroupKey)

	assert.Error(t, err)
}
