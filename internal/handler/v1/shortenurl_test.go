package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/repository"
	mockShortenURL "github.com/khaivutri/bookmark-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockSvc       func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL
		setupTestRequest   func(ctx *gin.Context)
		expectedStatusCode int
		expectedResponse   string
		expectedMessage    string
	}{
		{
			name: "normal case",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL {
				mockSvc := mockShortenURL.NewShortenURL(t)
				mockSvc.On("CreateCodeFromLink", ctx, "https://example.com", int64(3600)).Return("abc1234", nil).Once()
				return mockSvc
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"url":"https://example.com","exp":3600}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"code":"abc1234","message":"Shorten URL generated successfully!"}`,
		},
		{
			name: "returns error when input is invalid",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL {
				return mockShortenURL.NewShortenURL(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"url":"invalid-url","exp":4}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid input",
		},
		{
			name: "returns error when service fails",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL {
				mockSvc := mockShortenURL.NewShortenURL(t)
				mockSvc.On("CreateCodeFromLink", ctx, "https://example.com", int64(3600)).Return("", errors.New("service error")).Once()
				return mockSvc
			},

			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"url":"https://example.com","exp":3600}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Processing error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupTestRequest(ctx)

			mockSvc := tc.setupMockSvc(ctx, t)
			testHandler := NewShortenURL(mockSvc)

			testHandler.CreateShortenLink(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			if tc.expectedMessage != "" {
				var body struct {
					Message string `json:"message"`
				}
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
				assert.Equal(t, tc.expectedMessage, body.Message)
			} else {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}

func TestShortenURL_Redirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name 				string

		setupMockSvc     	func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL
		setupTestRequest 	func(ctx *gin.Context)

		expectedStatusCode 	int
		expectedLocation   	string
		expectedResponse   	string
	}{
		{
			name: "normal case",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL {
				mockSvc := mockShortenURL.NewShortenURL(t)
				mockSvc.On("GetLinkFromCode", ctx, "abc1234").Return("https://example.com", nil).Once()
				return mockSvc
			},

			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/abc1234", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "abc1234"}}
			},

			expectedStatusCode: http.StatusFound,
			expectedLocation:   "https://example.com",
		},
		{
			name: "returns error when code is empty",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL {
				return mockShortenURL.NewShortenURL(t)
			},

			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/", nil)
				ctx.Params = gin.Params{{Key: "code", Value: ""}}
			},

			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid input"}`,
		},
		{
			name: "returns error when code not found",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL {
				mockSvc := mockShortenURL.NewShortenURL(t)
				mockSvc.On("GetLinkFromCode", ctx, "abc1234").Return("", repository.ErrCodeNotFound).Once()
				return mockSvc
			},

			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/abc1234", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "abc1234"}}
			},

			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `{"error":"Code not found"}`,
		},
		{
			name: "returns error when service fails",

			setupMockSvc: func(ctx context.Context, t *testing.T) *mockShortenURL.ShortenURL {
				mockSvc := mockShortenURL.NewShortenURL(t)
				mockSvc.On("GetLinkFromCode", ctx, "abc1234").Return("", errors.New("service error")).Once()
				return mockSvc
			},

			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/abc1234", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "abc1234"}}
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Processing error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupTestRequest(ctx)

			mockSvc := tc.setupMockSvc(ctx, t)

			testHandler := NewShortenURL(mockSvc)

			testHandler.Redirect(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedLocation, rec.Header().Get("Location"))
			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}