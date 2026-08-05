package integrationtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/khaivutri/bookmark-service/internal/model"
	fixture "github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	jwtmocks "github.com/khaivutri/bookmark-service/pkg/jwtutils/mocks"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/khaivutri/bookmark-service/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	getBookmarksEndpoint = "/v1/bookmarks"
	getBookmarksUserID   = "b649b57b-b7b6-44e4-a233-74147ecf56ee"
)

type getBookmarksTestCase struct {
	name  string
	query string

	setupDB func(t *testing.T) *gorm.DB

	expectedStatus  int
	expectedPage    int
	expectedLimit   int
	expectedTotal   int64
	expectedCodes   []string
	expectedMessage string
}

func TestBookmarkEndpoint_GetBookmarks(t *testing.T) {
	testCases := []getBookmarksTestCase{
		{
			name:           "200 - returns the first page of bookmarks",
			query:          "?page=1&limit=1",
			setupDB:        setupGetBookmarksDB,
			expectedStatus: http.StatusOK,
			expectedPage:   1,
			expectedLimit:  1,
			expectedTotal:  2,
			expectedCodes:  []string{"code1"},
		},
		{
			name:           "200 - returns the second page of bookmarks",
			query:          "?page=2&limit=1",
			setupDB:        setupGetBookmarksDB,
			expectedStatus: http.StatusOK,
			expectedPage:   2,
			expectedLimit:  1,
			expectedTotal:  2,
			expectedCodes:  []string{"code2"},
		},
		{
			name:            "400 - rejects invalid pagination",
			query:           "?page=0&limit=10",
			setupDB:         setupGetBookmarksDB,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input",
		},
		{
			name:            "500 - returns processing error when database is down",
			query:           "?page=1&limit=10",
			setupDB:         getBookmarksDBDown,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := tc.setupDB(t)
			recorder := setupTestHTTP(t, db, http.MethodGet, getBookmarksEndpoint+tc.query, "")

			require.Equal(t, tc.expectedStatus, recorder.Code)
			assertGetBookmarksResponse(t, recorder, tc)
		})
	}
}

func TestBookmarkEndpoint_GetBookmarks_ReturnsCachedResponse(t *testing.T) {
	db := setupGetBookmarksDB(t)

	t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
	t.Setenv("SERVICE_NAME", "bookmark_service_test")
	require.NoError(t, validation.RegisterValidation())

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	jwtValidator := jwtmocks.NewJWTValidator(t)
	jwtValidator.
		On("ValidateJWT", getBookmarksToken).
		Return(jwt.MapClaims{"sub": getBookmarksUserID}, nil).
		Twice()

	testAPI := api.NewEngine(api.EngineOpts{
		App:    gin.New(),
		Cfg:    cfg,
		Redis:  redisPkg.InitMockRedis(t),
		DB:     db,
		JWTVal: jwtValidator,
	})

	firstResponse := serveGetBookmarksRequest(t, testAPI)
	require.Equal(t, http.StatusOK, firstResponse.Code)

	closeDatabase(t, db)

	secondResponse := serveGetBookmarksRequest(t, testAPI)
	require.Equal(t, http.StatusOK, secondResponse.Code)
	assertGetBookmarksResponse(t, secondResponse, getBookmarksTestCase{
		expectedPage:  1,
		expectedLimit: 1,
		expectedTotal: 2,
		expectedCodes: []string{"code1"},
	})
}

const getBookmarksToken = "get-bookmarks-test-token"

func serveGetBookmarksRequest(t *testing.T, testAPI http.Handler) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, getBookmarksEndpoint+"?page=1&limit=1", nil)
	req.Header.Set("Authorization", "Bearer "+getBookmarksToken)
	recorder := httptest.NewRecorder()
	testAPI.ServeHTTP(recorder, req)
	return recorder
}

func setupGetBookmarksDB(t *testing.T) *gorm.DB {
	t.Helper()
	return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
}

func getBookmarksDBDown(t *testing.T) *gorm.DB {
	t.Helper()
	db := setupGetBookmarksDB(t)
	closeDatabase(t, db)
	return db
}

func assertGetBookmarksResponse(t *testing.T, recorder *httptest.ResponseRecorder, tc getBookmarksTestCase) {
	t.Helper()
	if tc.expectedMessage != "" {
		assert.Contains(t, recorder.Body.String(), tc.expectedMessage)
		return
	}

	var response struct {
		Data       []model.Bookmark `json:"data"`
		Pagination struct {
			Page  int   `json:"page"`
			Limit int   `json:"limit"`
			Total int64 `json:"total"`
		} `json:"pagination"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	assert.Equal(t, tc.expectedPage, response.Pagination.Page)
	assert.Equal(t, tc.expectedLimit, response.Pagination.Limit)
	assert.Equal(t, tc.expectedTotal, response.Pagination.Total)
	assert.Equal(t, tc.expectedCodes, bookmarkCodes(response.Data))
}

func bookmarkCodes(bookmarks []model.Bookmark) []string {
	codes := make([]string, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		codes = append(codes, bookmark.Code)
	}
	return codes
}
