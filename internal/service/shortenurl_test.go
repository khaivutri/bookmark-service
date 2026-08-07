package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/khaivutri/bookmark-service/internal/repository"
	repoMocks "github.com/khaivutri/bookmark-service/internal/repository/mocks"
	utilsMocks "github.com/khaivutri/bookmark-service/pkg/utils/mocks"
	"github.com/stretchr/testify/require"
)

var (
	errGenFailed       = errors.New("gen error")
	errRedisUnexpected = errors.New("unexpected redis error")
)

const (
	testURL  = "https://example.com"
	testCode = "abc1234"
	testExp  = int64(3600)
)

type createCodeTestCase struct {
	name         string
	setupMocks   func(context.Context, *utilsMocks.GenCode, *repoMocks.URLStorage)
	expectedCode string
	expectedErr  error
}

func TestShortenURL_CreateCodeFromLink(t *testing.T) {
	tests := []createCodeTestCase{
		{
			name: "success on first attempt",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
				store.On("StoreURL", ctx, testCode, testURL, expiry()).Return(nil).Once()
			},
			expectedCode: testCode,
		},
		{
			name: "returns error when generator fails",
			setupMocks: func(_ context.Context, gen *utilsMocks.GenCode, _ *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return("", errGenFailed).Once()
			},
			expectedErr: errGenFailed,
		},
		{
			name: "returns error when storage GetURL fails unexpectedly",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", errRedisUnexpected).Once()
			},
			expectedErr: errRedisUnexpected,
		},
		{
			name: "retries generating code when collision occurs",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return("dup0001", nil).Once()
				store.On("GetURL", ctx, "dup0001").Return("https://old.example.com", nil).Once()
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
				store.On("StoreURL", ctx, testCode, testURL, expiry()).Return(nil).Once()
			},
			expectedCode: testCode,
		},
		{
			name: "returns error when storage StoreURL fails",
			setupMocks: func(ctx context.Context, gen *utilsMocks.GenCode, store *repoMocks.URLStorage) {
				gen.On("Generate", CODE_LEN).Return(testCode, nil).Once()
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
				store.On("StoreURL", ctx, testCode, testURL, expiry()).Return(errRedisUnexpected).Once()
			},
			expectedErr: errRedisUnexpected,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			genMock := utilsMocks.NewGenCode(t)
			storeMock := repoMocks.NewURLStorage(t)

			tc.setupMocks(ctx, genMock, storeMock)
			gotCode, err := NewURLStorage(storeMock, genMock).CreateCodeFromLink(ctx, testURL, testExp)

			requireCodeResult(t, gotCode, err, tc.expectedCode, tc.expectedErr)
		})
	}
}

func TestShortenURL_GetLinkFromCode(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(context.Context, *repoMocks.URLStorage)
		expectedURL string
		expectedErr error
	}{
		{
			name: "success",
			setupMocks: func(ctx context.Context, store *repoMocks.URLStorage) {
				store.On("GetURL", ctx, testCode).Return(testURL, nil).Once()
			},
			expectedURL: testURL,
		},
		{
			name: "returns error when code not found",
			setupMocks: func(ctx context.Context, store *repoMocks.URLStorage) {
				store.On("GetURL", ctx, testCode).Return("", repository.ErrCodeNotFound).Once()
			},
			expectedErr: repository.ErrCodeNotFound,
		},
		{
			name: "returns error when storage fails unexpectedly",
			setupMocks: func(ctx context.Context, store *repoMocks.URLStorage) {
				store.On("GetURL", ctx, testCode).Return("", errRedisUnexpected).Once()
			},
			expectedErr: errRedisUnexpected,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			storeMock := repoMocks.NewURLStorage(t)
			genMock := utilsMocks.NewGenCode(t)

			tc.setupMocks(ctx, storeMock)
			gotURL, err := NewURLStorage(storeMock, genMock).GetLinkFromCode(ctx, testCode)

			requireCodeResult(t, gotURL, err, tc.expectedURL, tc.expectedErr)
		})
	}
}

func expiry() time.Duration {
	return time.Duration(testExp) * time.Second
}

func requireCodeResult(t *testing.T, got string, err error, expected string, expectedErr error) {
	t.Helper()
	if expectedErr == nil {
		require.NoError(t, err)
	} else {
		require.ErrorIs(t, err, expectedErr)
	}
	require.Equal(t, expected, got)
}
