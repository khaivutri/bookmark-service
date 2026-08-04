package integrationtest

import (
	"bytes"
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
	addBookmarkEndpoint = "/v1/bookmarks"
	addBookmarkToken    = "add-bookmark-test-token"
	addBookmarkUserID   = "b649b57b-b7b6-44e4-a233-74147ecf56ee"
	addBookmarkBody     = `{"description":"Go documentation","url":"https://go.dev/doc/"}`
)

type addBookmarkTestCase struct {
	name string

	method       string
	body         string
	withAuth     bool
	databaseDown bool

	expectedStatus  int
	expectedMessage string
	expectedError   string
}

func TestBookmarkEndpoint_AddBookmark(t *testing.T) {
	testCases := []addBookmarkTestCase{
		{
			name:           "200 - creates bookmark successfully",
			method:         http.MethodPost,
			body:           addBookmarkBody,
			withAuth:       true,
			expectedStatus: http.StatusOK,
		},
		{
			name:            "400 - rejects invalid URL",
			method:          http.MethodPost,
			body:            `{"description":"Go documentation","url":"not-a-url"}`,
			withAuth:        true,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input",
		},
		{
			name:            "400 - rejects malformed JSON",
			method:          http.MethodPost,
			body:            `{invalid-json}`,
			withAuth:        true,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid input",
		},
		{
			name:           "401 - rejects request without authentication",
			method:         http.MethodPost,
			body:           addBookmarkBody,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
		{
			name:            "500 - returns processing error when database is down",
			method:          http.MethodPost,
			body:            addBookmarkBody,
			withAuth:        true,
			databaseDown:    true,
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupAddBookmarkDB(t)
			if tc.databaseDown {
				closeDatabase(t, db)
			}

			recorder := serveAddBookmarkRequest(t, db, tc)
			require.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedMessage != "" {
				assertResponseMessage(t, recorder, tc.expectedMessage)
				return
			}
			if tc.expectedError != "" {
				assertResponseError(t, recorder, tc.expectedError)
				return
			}

			assertCreatedBookmark(t, db, recorder)
		})
	}
}

func setupAddBookmarkDB(t *testing.T) *gorm.DB {
	t.Helper()
	return fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
}

func serveAddBookmarkRequest(t *testing.T, db *gorm.DB, tc addBookmarkTestCase) *httptest.ResponseRecorder {
	t.Helper()

	t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
	t.Setenv("SERVICE_NAME", "bookmark_service_test")
	require.NoError(t, validation.RegisterValidation())

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	jwtValidator := jwtmocks.NewJWTValidator(t)
	if tc.withAuth {
		jwtValidator.
			On("ValidateJWT", addBookmarkToken).
			Return(jwt.MapClaims{"sub": addBookmarkUserID}, nil).
			Once()
	}

	testAPI := api.NewEngine(api.EngineOpts{
		App:    gin.New(),
		Cfg:    cfg,
		Redis:  redisPkg.InitMockRedis(t),
		DB:     db,
		JWTVal: jwtValidator,
	})

	req := httptest.NewRequest(tc.method, addBookmarkEndpoint, bytes.NewBufferString(tc.body))
	if tc.withAuth {
		req.Header.Set("Authorization", "Bearer "+addBookmarkToken)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testAPI.ServeHTTP(recorder, req)
	return recorder
}

func assertCreatedBookmark(t *testing.T, db *gorm.DB, recorder *httptest.ResponseRecorder) {
	t.Helper()

	var response model.Bookmark
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.Code)
	assert.Equal(t, "Go documentation", response.Description)
	assert.Equal(t, "https://go.dev/doc/", response.URL)
	assert.Equal(t, addBookmarkUserID, response.UserID)

	var stored model.Bookmark
	require.NoError(t, db.First(&stored, "id = ?", response.ID).Error)
	assert.Equal(t, response.Code, stored.Code)
	assert.Equal(t, response.Description, stored.Description)
	assert.Equal(t, response.URL, stored.URL)
	assert.Equal(t, response.UserID, stored.UserID)
}

func assertResponseMessage(t *testing.T, recorder *httptest.ResponseRecorder, expected string) {
	t.Helper()

	var response struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	assert.Equal(t, expected, response.Message)
}

func assertResponseError(t *testing.T, recorder *httptest.ResponseRecorder, expected string) {
	t.Helper()

	var response struct {
		Error string `json:"error"`
	}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	assert.Equal(t, expected, response.Error)
}

func closeDatabase(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
}
