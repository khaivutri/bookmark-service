package user

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/model"
	mockSvc "github.com/khaivutri/bookmark-service/internal/service/user/mocks"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
)

const fakeUserID = "7d80f755-7dce-4c95-b8bf-75bb8e240ef2"

func setFakeUserID(ctx *gin.Context) {
	ctx.Set("claims", jwt.MapClaims{"sub": fakeUserID})
}

func TestHandler_GetSelfInfo(t *testing.T) {
	t.Parallel()

	stubUser := &model.User{
		Base:        fixture.GetTestBase(fakeUserID),
		UserName:    "johndoe",
		DisplayName: "John Doe",
		Email:       "john.doe@example.com",
	}

	testCases := []struct {
		name             string
		setupMockSvc     func(ctx context.Context, t *testing.T) *mockSvc.Service
		setupTestRequest func(ctx *gin.Context)
		expectedStatus   int
		expectedBody     string
	}{
		{
			name: "200 - gets current user successfully",
			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On("GetSelfInfo", ctx, fakeUserID).Return(stubUser, nil).Once()
				return svc
			},
			setupTestRequest: func(ctx *gin.Context) {
				setFakeUserID(ctx)
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"7d80f755-7dce-4c95-b8bf-75bb8e240ef2","display_name":"John Doe","username":"johndoe","email":"john.doe@example.com","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}`,
		},
		{
			name: "500 - service returns unexpected error",
			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On("GetSelfInfo", ctx, fakeUserID).Return(nil, errors.New("database connection lost")).Once()
				return svc
			},
			setupTestRequest: func(ctx *gin.Context) {
				setFakeUserID(ctx)
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Processing error"}`,
		},
		{
			name: "401 - missing claims",
			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid token"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			tc.setupTestRequest(ctx)

			handler := NewHandler(tc.setupMockSvc(ctx, t))
			handler.GetSelfInfo(ctx)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}

func TestHandler_UpdateSelfInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupMockSvc     func(ctx context.Context, t *testing.T) *mockSvc.Service
		setupTestRequest func(ctx *gin.Context)
		expectedStatus   int
		expectedBody     string
	}{
		{
			name: "200 - updates current user successfully",
			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On("UpdateSelfInfo", ctx, fakeUserID, "Jane Doe", "jane.doe@example.com").Return(nil).Once()
				return svc
			},
			setupTestRequest: func(ctx *gin.Context) {
				setFakeUserID(ctx)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBufferString(`{"display_name":"Jane Doe","email":"jane.doe@example.com"}`))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Edit current user successfully!"}`,
		},
		{
			name: "400 - invalid email",
			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				setFakeUserID(ctx)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBufferString(`{"display_name":"Jane Doe","email":"not-an-email"}`))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"Invalid input"`,
		},
		{
			name: "500 - service returns unexpected error",
			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On("UpdateSelfInfo", ctx, fakeUserID, "Jane Doe", "jane.doe@example.com").Return(errors.New("database connection lost")).Once()
				return svc
			},
			setupTestRequest: func(ctx *gin.Context) {
				setFakeUserID(ctx)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBufferString(`{"display_name":"Jane Doe","email":"jane.doe@example.com"}`))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Processing error"}`,
		},
		{
			name: "401 - missing claims",
			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBufferString(`{"display_name":"Jane Doe","email":"jane.doe@example.com"}`))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid token"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			tc.setupTestRequest(ctx)

			handler := NewHandler(tc.setupMockSvc(ctx, t))
			handler.UpdateSelfInfo(ctx)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.name == "400 - invalid email" {
				assert.True(t, strings.Contains(recorder.Body.String(), tc.expectedBody))
				return
			}
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
