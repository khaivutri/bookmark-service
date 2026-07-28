package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
	selfInfoEndpoint = "/v1/self/info"
	selfToken        = "self-test-token"
	selfUserID       = "b649b57b-b7b6-44e4-a233-74147ecf56ee"
)

func setupSelfDB(t *testing.T) *gorm.DB {
	t.Helper()
	return fixture.NewFixture(t, &fixture.UserCommonTest{})
}

// setupTestHTTP creates an authenticated request against the real HTTP stack.
func setupTestHTTP(t *testing.T, db *gorm.DB, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
	t.Setenv("SERVICE_NAME", "bookmark_service_test")
	require.NoError(t, validation.RegisterValidation())

	cfg, err := api.NewConfig()
	require.NoError(t, err)

	mockJWTValidator := jwtmocks.NewJWTValidator(t)
	if method != http.MethodPost {
		mockJWTValidator.
			On("ValidateJWT", selfToken).
			Return(jwt.MapClaims{"sub": selfUserID}, nil).
			Once()
	}

	testAPI := api.NewEngine(api.EngineOpts{
		App:    gin.New(),
		Cfg:    cfg,
		Redis:  redisPkg.InitMockRedis(t),
		DB:     db,
		JWTVal: mockJWTValidator,
	})

	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+selfToken)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	recorder := httptest.NewRecorder()
	testAPI.ServeHTTP(recorder, req)
	return recorder
}

func TestSelfEndPoint_GetSelfInfo(t *testing.T) {
	testCases := []struct {
		name               string
		method             string
		path               string
		setupDB            func(t *testing.T) *gorm.DB
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:               "200 - gets current user successfully",
			method:             http.MethodGet,
			path:               selfInfoEndpoint,
			setupDB:            setupSelfDB,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:   "500 - returns processing error when database is down",
			method: http.MethodGet,
			path:   selfInfoEndpoint,
			setupDB: func(t *testing.T) *gorm.DB {
				db := setupSelfDB(t)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
				return db
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "Processing error",
		},
		{
			name:               "405 - method not allowed",
			method:             http.MethodPost,
			path:               selfInfoEndpoint,
			setupDB:            setupSelfDB,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSqlDB := tc.setupDB(t)
			recorder := setupTestHTTP(t, mockSqlDB, tc.method, tc.path, "")

			require.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedMessage != "" {
				var response struct {
					Message string `json:"message"`
				}
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
				assert.Equal(t, tc.expectedMessage, response.Message)
				return
			}

			if tc.expectedStatusCode == http.StatusOK {
				var response model.User
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
				assert.Equal(t, selfUserID, response.ID)
				assert.Equal(t, "Test1", response.DisplayName)
				assert.Equal(t, "test1", response.UserName)
				assert.Equal(t, "test1@example.com", response.Email)
			}
		})
	}
}

func TestSelfEndPoint_UpdateSelfInfo(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		setupDB            func(t *testing.T) *gorm.DB
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:               "200 - updates current user successfully",
			body:               `{"display_name":"Jane Doe","email":"jane.doe@example.com"}`,
			setupDB:            setupSelfDB,
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Edit current user successfully!",
		},
		{
			name:               "400 - invalid input",
			body:               `{"display_name":"Jane Doe","email":"not-an-email"}`,
			setupDB:            setupSelfDB,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "500 - returns processing error when database is down",
			body: `{"display_name":"Jane Doe","email":"jane.doe@example.com"}`,
			setupDB: func(t *testing.T) *gorm.DB {
				db := setupSelfDB(t)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
				return db
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "Processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSqlDB := tc.setupDB(t)
			recorder := setupTestHTTP(t, mockSqlDB, http.MethodPut, selfInfoEndpoint, tc.body)

			require.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedStatusCode == http.StatusBadRequest {
				assert.True(t, strings.Contains(recorder.Body.String(), tc.expectedMessage))
			} else {
				var response struct {
					Message string `json:"message"`
				}
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
				assert.Equal(t, tc.expectedMessage, response.Message)
			}

			if tc.expectedStatusCode == http.StatusOK {
				var user model.User
				require.NoError(t, mockSqlDB.Where("id = ?", selfUserID).First(&user).Error)
				assert.Equal(t, "Jane Doe", user.DisplayName)
				assert.Equal(t, "jane.doe@example.com", user.Email)
			}
		})
	}
}
