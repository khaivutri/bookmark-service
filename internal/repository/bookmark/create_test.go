package bookmark

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name 			string
		setupDB 		func(t *testing.T) *gorm.DB
		inputBookmark 	*model.Bookmark

		expectedError 	error
		verifyFunc 		func(db *gorm.DB, bookmark *model.Bookmark)
	}{
		{
			name: "normal case",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmark: &model.Bookmark{
				Base: fixture.GetTestBase("11111111-1111-1111-1111-111111111111"),
				
				Description: "Test bookmark",
				URL:         "https://example.com",
				Code:        "test-bookmark",
				UserID:      "user1", 
			},
			expectedError: nil,
			verifyFunc: func(db *gorm.DB, expectedBookmark *model.Bookmark ) {
				var createdBookmark *model.Bookmark
				err := db.Where("code = ?", expectedBookmark.Code).First(&createdBookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, expectedBookmark.ID, createdBookmark.ID)
				assert.Equal(t, expectedBookmark.Description, createdBookmark.Description)
				assert.Equal(t, expectedBookmark.Code, createdBookmark.Code)
				assert.Equal(t, expectedBookmark.URL, createdBookmark.URL)
				assert.Equal(t, expectedBookmark.UserID, createdBookmark.UserID)
			},
		},
		{
			name: "duplicate code error",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			},
			inputBookmark: &model.Bookmark{
				Base: fixture.GetTestBase("22222222-2222-2222-2222-222222222222"),
				Description: "Test bookmark",
				URL:         "https://example.com",
				Code:        "code1", // duplicate code in fixture
				UserID:      "user1", 
			},
			expectedError: dbutils.ErrDuplicateBookmarkCode,
			verifyFunc: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewRepository(db)

			bookmark, err := repo.CreateBookmark(ctx, tc.inputBookmark)
			assert.ErrorIs(t, err, tc.expectedError)

			if tc.expectedError == nil {
				assert.NotNil(t, bookmark)
				if tc.verifyFunc != nil {
					tc.verifyFunc(db, tc.inputBookmark)
				}
			} else {
				assert.Nil(t, bookmark)
			}
		})
	}
}