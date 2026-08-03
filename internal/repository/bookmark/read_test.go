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

var (
	existingBookmark1 = &model.Bookmark{
		Base:        	fixture.GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ee"),
		Description: 	"description1",
		URL: 			"https://example.com",
		Code:		 	"code1",
		UserID: 		"b649b57b-b7b6-44e4-a233-74147ecf56ee",
	}
	existingBookmark2 = &model.Bookmark{
		Base: 			fixture.GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ef"),
		Description: 	"description2",
		URL: 			"https://example.com",
		Code: 			"code2",
		UserID: 		"b649b57b-b7b6-44e4-a233-74147ecf56ee",
	}
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
		t.Errorf("User ID mismatch: want %q, got %q", want.Code, got.Code)
	}
}

type bookmarkLookupCase struct {
	name 				string 
	code 				string 

	expectedBookmark 	*model.Bookmark
	expectedErr      	error
}

func testBookmarkLookup(t *testing.T, lookup func(*bookmarkRepo, context.Context, string) (*model.Bookmark, error), cases []bookmarkLookupCase) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, _ := newTestRepository(t)
			got, err := lookup(repo, t.Context(), tc.code)

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
			name:         		"existing bookmark - returns bookmark successfully",
			code:         		"code1",
			expectedBookmark: 	existingBookmark1,
			expectedErr:  		nil,
		},
		{
			name:         		"another existing bookmark - returns bookmark successfully",
			code:        		"code2",
			expectedBookmark: 	existingBookmark2,
			expectedErr:  		nil,
		},
		{
			name:         		"bookmark not found - returns parsed not found error",
			code:        		"non_existent_user",
			expectedBookmark: 	nil,
			expectedErr:  		dbutils.ErrRecordNotFound,
		},
		{
			name:         		"empty bookmark code - returns not found error",
			code:       		"",
			expectedBookmark: 	nil,
			expectedErr:  		dbutils.ErrRecordNotFound,
		},
	})
}

