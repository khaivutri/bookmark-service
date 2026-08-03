package bookmark

import (
	"errors"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	repoMocks "github.com/khaivutri/bookmark-service/internal/repository/bookmark/mocks"
	"github.com/stretchr/testify/require"
)

type getBookmarksTestCase struct {
	name          string
	userID        string
	page          int
	limit         int
	bookmarks     []*model.Bookmark
	total         int64
	repositoryErr error
}

func TestBookmarkService_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []getBookmarksTestCase{
		{
			name:      "returns bookmarks and total",
			userID:    "user-1",
			page:      2,
			limit:     2,
			bookmarks: []*model.Bookmark{bookmarkWithCode("code-1")},
			total:     3,
		},
		{
			name:          "returns repository error",
			userID:        "user-1",
			page:          1,
			limit:         10,
			repositoryErr: errors.New("get bookmarks failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			repo := repoMocks.NewBookmarkRepository(t)
			offset := (tc.page - 1) * tc.limit
			repo.On("GetBookmarks", ctx, tc.userID, tc.limit, offset).
				Return(tc.bookmarks, tc.total, tc.repositoryErr).Once()

			got, err := NewBookmarkService(repo, nil).GetBookmarks(ctx, tc.userID, tc.page, tc.limit)

			assertGetBookmarksResult(t, tc, got, err)
		})
	}
}

func assertGetBookmarksResult(t *testing.T, tc getBookmarksTestCase, got *GetBookmarkResponse, err error) {
	t.Helper()

	if tc.repositoryErr != nil {
		require.ErrorIs(t, err, tc.repositoryErr)
		require.Nil(t, got)
		return
	}

	require.NoError(t, err)
	require.Equal(t, &GetBookmarkResponse{Bookmarks: tc.bookmarks, Total: tc.total}, got)
}
