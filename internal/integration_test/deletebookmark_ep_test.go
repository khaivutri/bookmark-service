package integrationtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	deleteBookmarkEndpoint = "/v1/bookmarks"
	deleteBookmarkID       = "b649b57b-b7b6-44e4-a233-74147ecf56ee"
)

type deleteBookmarkEndpointTestCase struct {
	name           string
	path           string
	setupDB        func(t *testing.T) *gorm.DB
	expectedStatus int
	expectedBody   string
}

func TestBookmarkEndpoint_DeleteBookmark(t *testing.T) {
	testCases := []deleteBookmarkEndpointTestCase{
		{
			name:           "200 - deletes bookmark successfully",
			path:           deleteBookmarkEndpoint + "?id=" + deleteBookmarkID,
			setupDB:        setupDeleteBookmarkDB,
			expectedStatus: http.StatusOK,
			expectedBody:   "Success",
		},
		{
			name:           "400 - rejects missing bookmark ID",
			path:           deleteBookmarkEndpoint,
			setupDB:        setupDeleteBookmarkDB,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid input",
		},
		{
			name:           "404 - returns not found for missing bookmark",
			path:           deleteBookmarkEndpoint + "?id=missing-bookmark-id",
			setupDB:        setupDeleteBookmarkDB,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Bookmark not found",
		},
		{
			name:           "500 - returns processing error when database is down",
			path:           deleteBookmarkEndpoint + "?id=" + deleteBookmarkID,
			setupDB:        deleteBookmarkDBDown,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Fail to delete bookmark",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := tc.setupDB(t)
			recorder := setupTestHTTP(t, db, http.MethodDelete, tc.path, "")

			assertDeleteBookmarkEndpointResponse(t, db, recorder, tc)
		})
	}
}

func setupDeleteBookmarkDB(t *testing.T) *gorm.DB {
	t.Helper()
	return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
}

func deleteBookmarkDBDown(t *testing.T) *gorm.DB {
	t.Helper()
	db := setupDeleteBookmarkDB(t)
	closeDatabase(t, db)
	return db
}

func assertDeleteBookmarkEndpointResponse(
	t *testing.T,
	db *gorm.DB,
	recorder *httptest.ResponseRecorder,
	tc deleteBookmarkEndpointTestCase,
) {
	t.Helper()
	require.Equal(t, tc.expectedStatus, recorder.Code)
	assert.Contains(t, recorder.Body.String(), tc.expectedBody)

	if tc.expectedStatus != http.StatusOK {
		return
	}

	var deleted model.Bookmark
	err := db.First(&deleted, "id = ?", deleteBookmarkID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
