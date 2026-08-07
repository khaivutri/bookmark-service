package bookmark

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type bookmarkDeleteCase struct {
	name          string
	id            string
	setup         func(t *testing.T, db *gorm.DB)
	expectedError error
	expectDBError bool
}

func TestBookmarkRepo_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []bookmarkDeleteCase{
		{
			name: "deletes existing bookmark",
			id:   existingBookmark1.ID,
		},
		{
			name:          "returns not found for missing bookmark",
			id:            "missing-bookmark-id",
			expectedError: dbutils.ErrRecordNotFound,
		},
		{
			name: "returns database error when bookmarks table is missing",
			id:   existingBookmark1.ID,
			setup: func(t *testing.T, db *gorm.DB) {
				dropBookmarksTable(t, db)
			},
			expectDBError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, db := newTestRepository(t)
			if tc.setup != nil {
				tc.setup(t, db)
			}

			err := repo.DeleteBookmarkByID(t.Context(),"", tc.id)
			assertDeleteResult(t, db, tc, err)
		})
	}
}

func assertDeleteResult(t *testing.T, db *gorm.DB, tc bookmarkDeleteCase, err error) {
	t.Helper()

	if tc.expectedError != nil {
		assert.ErrorIs(t, err, tc.expectedError)
		return
	}
	if tc.expectDBError {
		assert.Error(t, err)
		return
	}

	require.NoError(t, err)
	var deleted model.Bookmark
	assert.ErrorIs(t, db.First(&deleted, "id = ?", tc.id).Error, gorm.ErrRecordNotFound)
}
