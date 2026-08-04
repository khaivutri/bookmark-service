package bookmark

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type bookmarkUpdateCase struct {
	name         string
	bookmarkID   string
	description  string
	url          string
	expectedErr  error
	verifyStored bool
}

func TestBookmarkRepo_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testBookmarkUpdates(t, []bookmarkUpdateCase{
		{
			name:         "updates existing bookmark successfully",
			bookmarkID:   existingBookmark1.ID,
			description:  "updated description",
			url:          "https://updated.example.com",
			verifyStored: true,
		},
		{
			name:        "non-existing id returns record not found",
			bookmarkID:  "c1111111-1111-1111-1111-111111111111",
			description: "some description",
			url:         "https://example.com",
			expectedErr: dbutils.ErrRecordNotFound,
		},
	})
}

func testBookmarkUpdates(t *testing.T, testCases []bookmarkUpdateCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, db := newTestRepository(t)
			err := repo.UpdateBookmark(t.Context(), tc.bookmarkID, tc.description, tc.url)
			assertUpdateResult(t, db, tc, err)
		})
	}
}

func assertUpdateResult(t *testing.T, db *gorm.DB, tc bookmarkUpdateCase, err error) {
	t.Helper()

	if tc.expectedErr != nil {
		assert.ErrorIs(t, err, tc.expectedErr)
		return
	}

	assert.NoError(t, err)
	if !tc.verifyStored {
		return
	}

	var got model.Bookmark
	assert.NoError(t, db.First(&got, "id = ?", tc.bookmarkID).Error)
	assert.Equal(t, tc.description, got.Description)
	assert.Equal(t, tc.url, got.URL)
}
