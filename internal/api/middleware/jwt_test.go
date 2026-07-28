package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/pkg/jwtutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputAuthHeader    string
		claims             jwt.MapClaims
		setupMockValidator func(t *testing.T, claims jwt.MapClaims) *mocks.JWTValidator

		expectedCode         int
		expectedResponseBody string
		expectedNextCalled   bool
	}{
		{
			name: "success",

			inputAuthHeader: "Bearer test-token-string",
			claims: jwt.MapClaims{
				"sub": "69c11072-9af6-4a5f-81d5-239d96154d5e",
			},
			setupMockValidator: func(t *testing.T, claims jwt.MapClaims) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateJWT", "test-token-string").Return(claims, nil).Once()
				return mockValidator
			},

			expectedCode:       http.StatusOK,
			expectedNextCalled: true,
		},
		{
			name: "returns unauthorized when authorization header is missing",

			setupMockValidator: func(t *testing.T, claims jwt.MapClaims) *mocks.JWTValidator {
				return mocks.NewJWTValidator(t)
			},

			expectedCode:         http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
		{
			name: "returns unauthorized when authorization scheme is not bearer",

			inputAuthHeader: "Basic test-token-string",
			setupMockValidator: func(t *testing.T, claims jwt.MapClaims) *mocks.JWTValidator {
				return mocks.NewJWTValidator(t)
			},

			expectedCode:         http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
		{
			name: "returns unauthorized when authorization header is malformed",

			inputAuthHeader: "Bearer",
			setupMockValidator: func(t *testing.T, claims jwt.MapClaims) *mocks.JWTValidator {
				return mocks.NewJWTValidator(t)
			},

			expectedCode:         http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
		{
			name: "returns unauthorized when jwt validation fails",

			inputAuthHeader: "Bearer invalid-token-string",
			setupMockValidator: func(t *testing.T, claims jwt.MapClaims) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateJWT", "invalid-token-string").Return(nil, errors.New("invalid token")).Once()
				return mockValidator
			},

			expectedCode:         http.StatusUnauthorized,
			expectedResponseBody: `{"error":"Unauthorized"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.inputAuthHeader != "" {
				req.Header.Set("Authorization", tc.inputAuthHeader)
			}

			mockValidator := tc.setupMockValidator(t, tc.claims)
			testMiddleware := NewJWTAuth(mockValidator)

			nextCalled := false
			router := gin.New()
			router.GET("/test", testMiddleware.JWTAuth(), func(ctx *gin.Context) {
				nextCalled = true

				claims, exists := ctx.Get("claims")
				assert.True(t, exists)
				assert.Equal(t, tc.claims, claims)

				ctx.Status(http.StatusOK)
			})

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expectedNextCalled, nextCalled)
			if tc.expectedResponseBody != "" {
				assert.JSONEq(t, tc.expectedResponseBody, rec.Body.String())
			}
		})
	}
}
