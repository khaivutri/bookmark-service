package integrationtest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/api"
	jwtmocks "github.com/khaivutri/bookmark-service/pkg/jwtutils/mocks"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/stretchr/testify/require"
)

func newMiddlewareTestConfig(t *testing.T) *api.Config {
	t.Helper()
	t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
	t.Setenv("SERVICE_NAME", "bookmark_service_test")
	cfg, err := api.NewConfig()
	require.NoError(t, err)
	return cfg
}

func TestJWTAuth_Integration(t *testing.T) {
	tests := []struct {
		name       string
		authorize  string
		setupMock  func(*jwtmocks.JWTValidator)
		wantStatus int
	}{
		{
			name:       "no token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token format",
			authorize:  "Basic token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "expired or invalid token",
			authorize: "Bearer invalid-token",
			setupMock: func(validator *jwtmocks.JWTValidator) {
				validator.On("ValidateJWT", "invalid-token").Return(jwt.MapClaims(nil), errors.New("token is invalid")).Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "valid token",
			authorize: "Bearer valid-token",
			setupMock: func(validator *jwtmocks.JWTValidator) {
				validator.On("ValidateJWT", "valid-token").Return(jwt.MapClaims{"sub": selfUserID}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupSelfDB(t)
			validator := jwtmocks.NewJWTValidator(t)
			if tt.setupMock != nil {
				tt.setupMock(validator)
			}

			engine := api.NewEngine(api.EngineOpts{
				App:    gin.New(),
				Cfg:    newMiddlewareTestConfig(t),
				Redis:  redisPkg.InitMockRedis(t),
				DB:     db,
				JWTVal: validator,
			})

			req := httptest.NewRequest(http.MethodGet, selfInfoEndpoint, nil)
			if tt.authorize != "" {
				req.Header.Set("Authorization", tt.authorize)
			}
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, req)

			require.Equal(t, tt.wantStatus, recorder.Code)
		})
	}
}

func TestRateLimiter_Integration(t *testing.T) {
	t.Run("limits public requests by IP and path", func(t *testing.T) {
		t.Setenv("RATE_LIMIT_IP_LIMIT", "3")
		engine := api.NewEngine(api.EngineOpts{
			App:   gin.New(),
			Cfg:   newMiddlewareTestConfig(t),
			Redis: redisPkg.InitMockRedis(t),
		})

		send := func(ip string) int {
			req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/somecode", nil)
			req.RemoteAddr = ip + ":1234"
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, req)
			return recorder.Code
		}

		require.Equal(t, http.StatusNotFound, send("192.168.1.1"))
		require.Equal(t, http.StatusNotFound, send("192.168.1.1"))
		require.Equal(t, http.StatusNotFound, send("192.168.1.1"))
		require.Equal(t, http.StatusTooManyRequests, send("192.168.1.1"))
		require.Equal(t, http.StatusNotFound, send("192.168.1.2"))
	})

	t.Run("limits protected requests by user", func(t *testing.T) {
		t.Setenv("RATE_LIMIT_USER_LIMIT", "2")
		validator := jwtmocks.NewJWTValidator(t)
		validator.On("ValidateJWT", "user-1-token").Return(jwt.MapClaims{"sub": selfUserID}, nil).Times(3)
		validator.On("ValidateJWT", "user-2-token").Return(jwt.MapClaims{"sub": "b649b57b-b7b6-44e4-a233-74147ecf56ef"}, nil).Once()

		engine := api.NewEngine(api.EngineOpts{
			App:    gin.New(),
			Cfg:    newMiddlewareTestConfig(t),
			Redis:  redisPkg.InitMockRedis(t),
			DB:     setupSelfDB(t),
			JWTVal: validator,
		})

		send := func(token string) int {
			req := httptest.NewRequest(http.MethodGet, selfInfoEndpoint, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, req)
			return recorder.Code
		}

		require.Equal(t, http.StatusOK, send("user-1-token"))
		require.Equal(t, http.StatusOK, send("user-1-token"))
		require.Equal(t, http.StatusTooManyRequests, send("user-1-token"))
		require.Equal(t, http.StatusOK, send("user-2-token"))
	})
}
