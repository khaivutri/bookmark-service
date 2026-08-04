package bookmark

import (
	"context"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newTestRepository(t *testing.T) (*bookmarkRepo, *gorm.DB) {
	t.Helper()

	db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})

	return &bookmarkRepo{db: db}, db
}

func assertBookmarkEqual(t *testing.T, want, got *model.Bookmark) {
	t.Helper()

	if got == nil {
		t.Fatalf("expected bookmark %+v, got nil", want)
	}

	if got.ID != want.ID {
		t.Errorf("ID mismatch: want %q, got %q", want.ID, got.ID)
	}
	if got.Description != want.Description {
		t.Errorf("Description mismatch: want %q, got %q", want.Description, got.Description)
	}
	if got.URL != want.URL {
		t.Errorf("DisplayName mismatch: want %q, got %q", want.URL, got.URL)
	}
	if got.Code != want.Code {
		t.Errorf("Code mismatch: want %q, got %q", want.Code, got.Code)
	}
	if got.UserID != want.UserID {
		t.Errorf("User ID mismatch: want %q, got %q", want.UserID, got.UserID)
	}
}

type bookmarkLookupCase struct {
	name             string
	value            string
	expectedBookmark *model.Bookmark
	expectedErr      error
}

func testBookmarkLookup(
	t *testing.T,
	lookup func(*bookmarkRepo, context.Context, string) (*model.Bookmark, error),
	cases []bookmarkLookupCase,
) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, _ := newTestRepository(t)
			got, err := lookup(repo, t.Context(), tc.value)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, got, "expected nil bookmark on error")
				return
			}

			assert.NoError(t, err)
			assertBookmarkEqual(t, tc.expectedBookmark, got)
		})
	}
}

func TestBookmarkRepo_GetBookmarkByCode(t *testing.T) {
	t.Parallel()

	testBookmarkLookup(t, func(repo *bookmarkRepo, ctx context.Context, code string) (*model.Bookmark, error) {
		return repo.GetBookmarkByCode(ctx, code)
	}, []bookmarkLookupCase{
		{
			name:             "existing bookmark - returns bookmark successfully",
			value:            "code1",
			expectedBookmark: existingBookmark1,
		},
		{
			name:             "another existing bookmark - returns bookmark successfully",
			value:            "code2",
			expectedBookmark: existingBookmark2,
		},
		{
			name:        "bookmark not found - returns parsed not found error",
			value:       "non_existent_user",
			expectedErr: dbutils.ErrRecordNotFound,
		},
		{
			name:        "empty bookmark code - returns not found error",
			value:       "",
			expectedErr: dbutils.ErrRecordNotFound,
		},
	})
}

func TestBookmarkRepo_GetBookmarkByID(t *testing.T) {
	t.Parallel()

	testBookmarkLookup(t, func(repo *bookmarkRepo, ctx context.Context, id string) (*model.Bookmark, error) {
		return repo.GetBookmarkByID(ctx, id)
	}, []bookmarkLookupCase{
		{
			name:             "existing bookmark - returns bookmark successfully",
			value:            existingBookmark1.ID,
			expectedBookmark: existingBookmark1,
		},
		{
			name:             "another existing bookmark - returns bookmark successfully",
			value:            existingBookmark2.ID,
			expectedBookmark: existingBookmark2,
		},
		{
			name:        "bookmark not found - returns parsed not found error",
			value:       "non-existent-bookmark-id",
			expectedErr: dbutils.ErrRecordNotFound,
		},
		{
			name:        "empty bookmark ID - returns not found error",
			value:       "",
			expectedErr: dbutils.ErrRecordNotFound,
		},
	})
}
