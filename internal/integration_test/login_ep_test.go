package integrationtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/khaivutri/bookmark-service/internal/model"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/khaivutri/bookmark-service/pkg/sqldb"
	"github.com/khaivutri/bookmark-service/pkg/utils"
	"github.com/khaivutri/bookmark-service/pkg/validation"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	loginEndpoint      = "/v1/users/login"
	loginUsername      = "johndoe"
	loginPassword      = "Password123@"
	loginGeneratedJWT  = "jwt-token-string"
	loginUserID        = "69c11072-9af6-4a5f-81d5-239d96154d5e"
	loginUserEmail     = "john.doe@example.com"
	loginUserNameInDB  = loginUsername
	loginDisplayNameDB = "John Doe"
)

type loginJWTGenerator struct {
	token string
	err   error
}

type loginResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (g loginJWTGenerator) GenerateJWT(jwt.MapClaims) (string, error) {
	return g.token, g.err
}

func TestLoginEndpoint(t *testing.T) {
	testCases := []struct {
		name string

		reqMethod string
		reqPath   string
		reqBody   string

		setupDB      func(t *testing.T) *gorm.DB
		jwtGenerator loginJWTGenerator

		expectedStatusCode int
		expectedMessage    string
		expectedToken      string
	}{
		{
			name:      "200 - logs in successfully",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{"username":"johndoe","password":"Password123@"}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Logged in successfully!",
			expectedToken:      loginGeneratedJWT,
		},
		{
			name:      "400 - invalid input when body is empty",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name:      "400 - invalid input when json body is malformed",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{invalid-json}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name:      "400 - invalid input when password violates policy",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{"username":"johndoe","password":"password123"}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name:      "400 - invalid credential when username does not exist",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{"username":"unknownuser","password":"Password123@"}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid credential",
		},
		{
			name:      "400 - invalid credential when password does not match",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{"username":"johndoe","password":"Wrong123@"}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid credential",
		},
		{
			name:      "500 - returns processing error when jwt generation fails",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{"username":"johndoe","password":"Password123@"}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				err: errors.New("sign token failed"),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "Processing error",
		},
		{
			name:      "500 - returns processing error when database is down",
			reqMethod: http.MethodPost,
			reqPath:   loginEndpoint,
			reqBody:   `{"username":"johndoe","password":"Password123@"}`,
			setupDB: func(t *testing.T) *gorm.DB {
				db := setupLoginDB(t)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
				return db
			},
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "Processing error",
		},
		{
			name:      "405 - method not allowed",
			reqMethod: http.MethodGet,
			reqPath:   loginEndpoint,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:      "404 - route not found",
			reqMethod: http.MethodPost,
			reqPath:   "/v1/users/login-wrong",
			reqBody:   `{"username":"johndoe","password":"Password123@"}`,
			setupDB:   setupLoginDB,
			jwtGenerator: loginJWTGenerator{
				token: loginGeneratedJWT,
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
			t.Setenv("SERVICE_NAME", "bookmark_service_test")

			require.NoError(t, validation.RegisterValidation())

			cfg, err := api.NewConfig()
			require.NoError(t, err)

			testAPI := api.NewEngine(api.EngineOpts{
				App:    gin.New(),
				Cfg:    cfg,
				Redis:  redisPkg.InitMockRedis(t),
				DB:     tc.setupDB(t),
				JWTGen: tc.jwtGenerator,
			})

			req := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			testAPI.ServeHTTP(recorder, req)

			require.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedMessage == "" && tc.expectedToken == "" {
				return
			}

			var response loginResponse
			require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
			require.Equal(t, tc.expectedMessage, response.Message)
			require.Equal(t, tc.expectedToken, response.Data)
		})
	}
}

func setupLoginDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := sqldb.InitMockDB(t)
	require.NoError(t, db.AutoMigrate(&model.User{}))

	hashedPassword, err := utils.NewHasher().Hash(loginPassword)
	require.NoError(t, err)

	user := &model.User{
		ID:          loginUserID,
		DisplayName: loginDisplayNameDB,
		UserName:    loginUserNameInDB,
		Password:    hashedPassword,
		Email:       loginUserEmail,
	}

	require.NoError(t, db.Session(&gorm.Session{SkipHooks: true}).Create(user).Error)

	return db
}
