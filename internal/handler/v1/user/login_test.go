package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	serviceUser "github.com/khaivutri/bookmark-service/internal/service/user"
	mockSvc "github.com/khaivutri/bookmark-service/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockSvc     func(ctx context.Context, t *testing.T) *mockSvc.Service
		setupTestRequest func(ctx *gin.Context)

		expectedStatusCode int
		expectedMessage    string
		expectedToken      string
	}{
		{
			name: "200 - login successfully",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On("Login", ctx, "johndoe", "Password123@").Return("jwt-token-string", nil).Once()
				return svc
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"username":"johndoe","password":"Password123@"}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Logged in successfully!",
			expectedToken:      "jwt-token-string",
		},
		{
			name: "400 - missing required fields",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - username is too short",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"username":"jo","password":"Password123@"}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - password does not satisfy policy",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"username":"johndoe","password":"password123"}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - malformed json body",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{invalid-json}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - invalid credential",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On("Login", ctx, "johndoe", "Password123@").Return("", serviceUser.ErrInvalidCredential).Once()
				return svc
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"username":"johndoe","password":"Password123@"}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid credential",
		},
		{
			name: "500 - service returns unexpected error",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On("Login", ctx, "johndoe", "Password123@").Return("", errors.New("database connection lost")).Once()
				return svc
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"username":"johndoe","password":"Password123@"}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "Processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupTestRequest(ctx)

			svc := tc.setupMockSvc(ctx, t)
			handler := NewHandler(svc)

			handler.Login(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			var respBody struct {
				Message string `json:"message"`
				Data    string `json:"data"`
			}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respBody))
			assert.Equal(t, tc.expectedMessage, respBody.Message)
			assert.Equal(t, tc.expectedToken, respBody.Data)
		})
	}
}
