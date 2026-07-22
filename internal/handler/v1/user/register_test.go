package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/model"
	mockSvc "github.com/khaivutri/bookmark-service/internal/service/user/mocks"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/khaivutri/bookmark-service/pkg/validation"
	"github.com/stretchr/testify/assert"
)
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
 
	if err := validation.RegisterValidation(); err != nil {
		panic("failed to register custom validators: " + err.Error())
	}
 
	os.Exit(m.Run())
}
 
func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	// stubUser is a reusable user model returned by the mock service on success.
	stubUser := &model.User{
		ID:          "7d80f755-7dce-4c95-b8bf-75bb8e240ef2",
		UserName:    "johndoe",
		DisplayName: "John Doe",
		Email:       "john.doe@example.com",
	}

	testCases := []struct {
		name 				string

		setupMockSvc     	func(ctx context.Context, t *testing.T) *mockSvc.Service
		setupTestRequest 	func(ctx *gin.Context)

		expectedStatusCode 	int
		expectedMessage    	string
		expectedResponse  	string 
	}{
		// ──────────────────────────────────────────────────────────────
		// Happy path
		// ──────────────────────────────────────────────────────────────
		{
			name: "201 – registers user successfully",

			setupMockSvc: func(ctx context.Context,t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On(
					"CreateUser",
					ctx,
					"johndoe", "John Doe", "Password123@", "john.doe@example.com",
				).Return(stubUser, nil).Once()
				return svc
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"johndoe","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusCreated,
			expectedMessage:    "User registered successfully!",
		},

		// ──────────────────────────────────────────────────────────────
		// Input validation – service must never be called on 400s
		// ──────────────────────────────────────────────────────────────
		{
			name: "400 – missing required fields (empty body)",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t) // no calls expected
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 – username too short (< 3 chars)",

			setupMockSvc: func(ctx context.Context ,t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"ab","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 – username too long (> 20 chars)",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"aaaaaaaaaaaaaaaaaaaaa","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 – invalid email format",

			setupMockSvc: func(ctx context.Context ,t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"johndoe","display_name":"John Doe","password":"Password123@","email":"not-an-email"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 – password missing special character",

			setupMockSvc: func(ctx context.Context ,t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"johndoe","display_name":"John Doe","password":"Password123","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 – password missing uppercase letter",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"johndoe","display_name":"John Doe","password":"password123@","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "400 – malformed JSON body",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				return mockSvc.NewService(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{invalid-json}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},

		// ──────────────────────────────────────────────────────────────
		// Conflict – duplicate username / email
		// ──────────────────────────────────────────────────────────────
		{
			name: "409 – username already exists",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On(
					"CreateUser",
					ctx,
					"johndoe", "John Doe", "Password123@", "john.doe@example.com",
				).Return(nil, dbutils.ErrDuplicateUserName).Once()
				return svc
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"johndoe","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusConflict,
			expectedMessage:    "username already exists",
		},
		{
			name: "409 – email already exists",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On(
					"CreateUser",
					ctx,
					"johndoe", "John Doe", "Password123@", "john.doe@example.com",
				).Return(nil, dbutils.ErrDuplicateEmail).Once()
				return svc
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"johndoe","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusConflict,
			expectedMessage:    "email already exists",
		},

		// ──────────────────────────────────────────────────────────────
		// Internal server error
		// ──────────────────────────────────────────────────────────────
		{
			name: "500 – service returns unexpected error",

			setupMockSvc: func(ctx context.Context,t *testing.T) *mockSvc.Service {
				svc := mockSvc.NewService(t)
				svc.On(
					"CreateUser",
					ctx,
					"johndoe", "John Doe", "Password123@", "john.doe@example.com",
				).Return(nil, errors.New("database connection lost")).Once()
				return svc
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(
					`{"username":"johndoe","display_name":"John Doe","password":"Password123@","email":"john.doe@example.com"}`,
				)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", body)
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

			handler.Register(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
				return
			}

			if tc.expectedMessage != "" {
				var respBody struct {
					Message string `json:"message"`
				}
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &respBody))
				assert.Equal(t, tc.expectedMessage, respBody.Message)
			}
		})
	}
}

