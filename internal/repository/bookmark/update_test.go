package bookmark

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type bookmarkUpdateCase struct {
	name         string
	input        *model.Bookmark
	expectedErr  error
	verifyStored bool
}

func TestBookmarkRepo_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testBookmarkUpdates(t, []bookmarkUpdateCase{
		{
			name: "updates existing bookmark successfully",
			input: &model.Bookmark{
				Base:        fixture.GetTestBase(existingBookmark1.ID),
				Description: "updated description",
				URL:         "https://updated.example.com",
				Code:        existingBookmark1.Code,
				UserID:      existingBookmark1.UserID,
			},
			verifyStored: true,
		},
		{
			name: "new id creates a bookmark",
			input: &model.Bookmark{
				Base:        fixture.GetTestBase("c1111111-1111-1111-1111-111111111111"),
				Description: "new description",
				URL:         "https://new.example.com",
				Code:        "new-code",
				UserID:      existingBookmark1.UserID,
			},
			verifyStored: true,
		},
		{
			name: "duplicate code returns parsed database error",
			input: &model.Bookmark{
				Base:        fixture.GetTestBase(existingBookmark1.ID),
				Description: existingBookmark1.Description,
				URL:         existingBookmark1.URL,
				Code:        existingBookmark2.Code,
				UserID:      existingBookmark1.UserID,
			},
			expectedErr: dbutils.ErrDuplicateBookmarkCode,
		},
	})
}

func testBookmarkUpdates(t *testing.T, testCases []bookmarkUpdateCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, db := newTestRepository(t)
			err := repo.UpdateBookmark(t.Context(), tc.input)
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
	assert.NoError(t, db.First(&got, "id = ?", tc.input.ID).Error)
	assertBookmarkEqual(t, tc.input, &got)
}
