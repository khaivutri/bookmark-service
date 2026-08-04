package bookmark

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type getBookmarksCase struct {
	name          string
	userID        string
	limit         int
	offset        int
	setup         func(t *testing.T, db *gorm.DB)
	expected      []*model.Bookmark
	expectedCount int64
	expectedError error
}

func TestBookmarkRepo_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []getBookmarksCase{
		{
			name:          "returns bookmarks ordered by creation time",
			userID:        existingBookmark1.UserID,
			limit:         2,
			expected:      []*model.Bookmark{existingBookmark1, existingBookmark2},
			expectedCount: 2,
		},
		{
			name:          "applies offset and limit",
			userID:        existingBookmark1.UserID,
			limit:         1,
			offset:        1,
			expected:      []*model.Bookmark{existingBookmark2},
			expectedCount: 2,
		},
		{
			name:          "returns empty result for user without bookmarks",
			userID:        "user-without-bookmarks",
			limit:         2,
			expected:      []*model.Bookmark{},
			expectedCount: 0,
		},
		{
			name:          "returns get error when bookmarks query fails",
			userID:        existingBookmark1.UserID,
			limit:         2,
			setup:         dropBookmarksTable,
			expectedError: ErrGetBookMarks,
		},
		{
			name:          "returns count error when count query fails",
			userID:        existingBookmark1.UserID,
			limit:         2,
			setup:         dropTableAfterBookmarksQuery,
			expectedError: ErrCount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, db := newTestRepository(t)
			if tc.setup != nil {
				tc.setup(t, db)
			}

			got, count, err := repo.GetBookmarks(t.Context(), tc.userID, tc.limit, tc.offset)

			assertGetBookmarksResult(t, tc, got, count, err)
		})
	}
}

func assertGetBookmarksResult(t *testing.T, tc getBookmarksCase, got []*model.Bookmark, count int64, err error) {
	t.Helper()

	if tc.expectedError != nil {
		assert.ErrorIs(t, err, tc.expectedError)
		assert.Nil(t, got)
		assert.Zero(t, count)
		return
	}

	require.NoError(t, err)
	assert.Equal(t, tc.expected, got)
	assert.Equal(t, tc.expectedCount, count)
}

func dropBookmarksTable(t *testing.T, db *gorm.DB) {
	t.Helper()
	require.NoError(t, db.Migrator().DropTable(&model.Bookmark{}))
}

func dropTableAfterBookmarksQuery(t *testing.T, db *gorm.DB) {
	t.Helper()

	dropped := false
	require.NoError(t, db.Callback().Query().After("gorm:query").Register("bookmark_test_drop_table", func(queryDB *gorm.DB) {
		if dropped || queryDB.Statement.Table != "bookmarks" {
			return
		}
		dropped = true
		require.NoError(t, queryDB.Session(&gorm.Session{NewDB: true}).Migrator().DropTable(&model.Bookmark{}))
	}))
	t.Cleanup(func() {
		db.Callback().Query().Remove("bookmark_test_drop_table")
	})
}
