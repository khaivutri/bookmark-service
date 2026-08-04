package integrationtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	fixture "github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	updateBookmarkEndpoint = "/v1/bookmarks"
	updateBookmarkID       = "b649b57b-b7b6-44e4-a233-74147ecf56ee"
	updateBookmarkQuery    = "?id=" + updateBookmarkID
	updateBookmarkBody     = `{"description":"Updated description","url":"https://go.dev/doc/"}`
)

type updateBookmarkEndpointTestCase struct {
	name            string
	body            string
	path            string
	setupDB         func(t *testing.T) *gorm.DB
	expectedStatus  int
	expectedMessage string
}

func TestBookmarkEndpoint_UpdateBookmark(t *testing.T) {
	testCases := []updateBookmarkEndpointTestCase{
		{
			name:            "200 - updates bookmark successfully",
			body:            updateBookmarkBody,
			path:            updateBookmarkEndpoint + updateBookmarkQuery,
			setupDB:         setupUpdateBookmarkDB,
			expectedStatus:  http.StatusOK,
			expectedMessage: "Success",
		},
		{
			name:            "400 - rejects invalid URL",
			body:            `{"description":"Updated description","url":"not-a-url"}`,
			path:            updateBookmarkEndpoint + updateBookmarkQuery,
			setupDB:         setupUpdateBookmarkDB,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input",
		},
		{
			name:            "400 - rejects missing bookmark ID",
			body:            updateBookmarkBody,
			path:            updateBookmarkEndpoint,
			setupDB:         setupUpdateBookmarkDB,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input",
		},
		{
			name:            "400 - rejects malformed JSON",
			body:            `{invalid-json}`,
			path:            updateBookmarkEndpoint + updateBookmarkQuery,
			setupDB:         setupUpdateBookmarkDB,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input",
		},
		{
			name:            "500 - returns processing error when database is down",
			body:            updateBookmarkBody,
			path:            updateBookmarkEndpoint + updateBookmarkQuery,
			setupDB:         updateBookmarkDBDown,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Fail to update bookmark",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := tc.setupDB(t)
			recorder := setupTestHTTP(t, db, http.MethodPut, tc.path, tc.body)

			assertUpdateBookmarkEndpointResponse(t, db, recorder, tc)
		})
	}
}

func setupUpdateBookmarkDB(t *testing.T) *gorm.DB {
	t.Helper()
	return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
}

func updateBookmarkDBDown(t *testing.T) *gorm.DB {
	t.Helper()
	db := setupUpdateBookmarkDB(t)
	closeDatabase(t, db)
	return db
}

func assertUpdateBookmarkEndpointResponse(
	t *testing.T,
	db *gorm.DB,
	recorder *httptest.ResponseRecorder,
	tc updateBookmarkEndpointTestCase,
) {
	t.Helper()
	require.Equal(t, tc.expectedStatus, recorder.Code)
	if tc.expectedStatus != http.StatusOK {
		assert.Contains(t, recorder.Body.String(), tc.expectedMessage)
		return
	}

	var response struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	assert.Equal(t, tc.expectedMessage, response.Message)

	var bookmark model.Bookmark
	require.NoError(t, db.First(&bookmark, "id = ?", updateBookmarkID).Error)
	assert.Equal(t, "Updated description", bookmark.Description)
	assert.Equal(t, "https://go.dev/doc/", bookmark.URL)
}
