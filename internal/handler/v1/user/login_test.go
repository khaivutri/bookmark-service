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

		setupMockSvc func(ctx context.Context, t *testing.T, svc *mockSvc.Service)
		requestBody  string

		expectedStatusCode int
		expectedMessage    string
		expectedToken      string
	}{
		{
			name: "200 - login successfully",

			setupMockSvc: func(ctx context.Context, t *testing.T, svc *mockSvc.Service) {
				svc.On("Login", ctx, "johndoe", "Password123@").Return("jwt-token-string", nil).Once()
			},
			requestBody: `{"username":"johndoe","password":"Password123@"}`,

			expectedStatusCode: http.StatusOK,
			expectedMessage:    "Logged in successfully!",
			expectedToken:      "jwt-token-string",
		},
		{
			name: "400 - missing required fields",

			setupMockSvc: func(context.Context, *testing.T, *mockSvc.Service) {},
			requestBody:  `{}`,

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - username is too short",

			setupMockSvc: func(context.Context, *testing.T, *mockSvc.Service) {},
			requestBody:  `{"username":"jo","password":"Password123@"}`,

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - password does not satisfy policy",

			setupMockSvc: func(context.Context, *testing.T, *mockSvc.Service) {},
			requestBody:  `{"username":"johndoe","password":"password123"}`,

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - malformed json body",

			setupMockSvc: func(context.Context, *testing.T, *mockSvc.Service) {},
			requestBody:  `{invalid-json}`,

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 - invalid credential",

			setupMockSvc: func(ctx context.Context, t *testing.T, svc *mockSvc.Service) {
				svc.On("Login", ctx, "johndoe", "Password123@").Return("", serviceUser.ErrInvalidCredential).Once()
			},
			requestBody: `{"username":"johndoe","password":"Password123@"}`,

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "invalid credential",
		},
		{
			name: "500 - service returns unexpected error",

			setupMockSvc: func(ctx context.Context, t *testing.T, svc *mockSvc.Service) {
				svc.On("Login", ctx, "johndoe", "Password123@").Return("", errors.New("database connection lost")).Once()
			},
			requestBody: `{"username":"johndoe","password":"Password123@"}`,

			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "Processing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBufferString(tc.requestBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			svc := mockSvc.NewService(t)
			tc.setupMockSvc(ctx, t, svc)
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
