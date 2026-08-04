package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	repoMocks "github.com/khaivutri/bookmark-service/internal/repository/bookmark/mocks"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	utilsMocks "github.com/khaivutri/bookmark-service/pkg/utils/mocks"
	"github.com/stretchr/testify/require"
)

var (
	errGenFailed      = errors.New("gen error")
	errRepoUnexpected = errors.New("unexpected repo error")
)

const (
	testDescription = "test description"
	testURL         = "https://example.com"
	testUserID      = "test-user-id"
)

func TestBookmarkService_AddBookmark(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		setupMocks  func(context.Context, *repoMocks.BookmarkRepository, *utilsMocks.GenCode)
		expected    *model.Bookmark
		expectedErr error
	}{
		{
			name: "success on first attempt",
			setupMocks: func(ctx context.Context, repo *repoMocks.BookmarkRepository, gen *utilsMocks.GenCode) {
				expectAvailableCode(ctx, repo, gen, "uniqueCode")
				expectCreateBookmark(ctx, repo, "uniqueCode", nil)
			},
			expected: bookmarkWithCode("uniqueCode"),
		},
		{
			name: "returns error when generator fails",
			setupMocks: func(_ context.Context, _ *repoMocks.BookmarkRepository, gen *utilsMocks.GenCode) {
				gen.On("Generate", 10).Return("", errGenFailed).Once()
			},
			expectedErr: errGenFailed,
		},
		{
			name: "returns error when repo CreateBookmark fails",
			setupMocks: func(ctx context.Context, repo *repoMocks.BookmarkRepository, gen *utilsMocks.GenCode) {
				expectAvailableCode(ctx, repo, gen, "uniqueCode")
				repo.On("CreateBookmark", ctx, bookmarkWithCode("uniqueCode")).Return(nil, errRepoUnexpected).Once()
			},
			expectedErr: errRepoUnexpected,
		},
		{
			name: "retries generating code when collision occurs",
			setupMocks: func(ctx context.Context, repo *repoMocks.BookmarkRepository, gen *utilsMocks.GenCode) {
				gen.On("Generate", 10).Return("collidedCode", nil).Once()
				repo.On("GetBookmarkByCode", ctx, "collidedCode").Return(bookmarkWithCode("collidedCode"), nil).Once()
				expectAvailableCode(ctx, repo, gen, "uniqueCode")
				expectCreateBookmark(ctx, repo, "uniqueCode", nil)
			},
			expected: bookmarkWithCode("uniqueCode"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := repoMocks.NewBookmarkRepository(t)
			gen := utilsMocks.NewGenCode(t)
			tc.setupMocks(ctx, repo, gen)

			result, err := NewBookmarkService(repo, gen).AddBookmark(ctx, testDescription, testURL, testUserID)

			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expected, result)
		})
	}
}

func bookmarkCreate() model.Bookmark {
	return model.Bookmark{
		Description: testDescription,
		URL:         testURL,
		UserID:      testUserID,
	}
}

func bookmarkWithCode(code string) *model.Bookmark {
	bookmark := bookmarkCreate()
	bookmark.Code = code
	return &bookmark
}

func expectAvailableCode(ctx context.Context, repo *repoMocks.BookmarkRepository, gen *utilsMocks.GenCode, code string) {
	gen.On("Generate", 10).Return(code, nil).Once()
	repo.On("GetBookmarkByCode", ctx, code).Return(nil, dbutils.ErrRecordNotFound).Once()
}

func expectCreateBookmark(ctx context.Context, repo *repoMocks.BookmarkRepository, code string, err error) {
	repo.On("CreateBookmark", ctx, bookmarkWithCode(code)).Return(bookmarkWithCode(code), err).Once()
}
