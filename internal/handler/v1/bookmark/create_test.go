package bookmark

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/service/bookmark/mocks"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkHandler_AddBookmark(t *testing.T) {
	t.Parallel()

	const (
		requestBody = `{"description":"Go documentation","url":"https://go.dev/doc/"}`
		userID      = "user-123"
	)

	createdBookmark := &model.Bookmark{
		Description: "Go documentation",
		URL:         "https://go.dev/doc/",
		Code:        "go-docs",
		UserID:      userID,
	}

	tests := []struct {
		name             string
		body             string
		withClaims       bool
		callService      bool
		serviceResult    *model.Bookmark
		serviceErr       error
		wantStatus       int
		wantResponseBody string
		wantResponsePart string
	}{
		{
			name:             "creates bookmark",
			body:             requestBody,
			withClaims:       true,
			callService:      true,
			serviceResult:    createdBookmark,
			wantStatus:       http.StatusOK,
			wantResponseBody: `{"id":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","description":"Go documentation","url":"https://go.dev/doc/","code":"go-docs","user_id":"user-123"}`,
		},
		{
			name:             "returns conflict when bookmark code is duplicated",
			body:             requestBody,
			withClaims:       true,
			callService:      true,
			serviceErr:       dbutils.ErrDuplicateBookmarkCode,
			wantStatus:       http.StatusConflict,
			wantResponseBody: `{"message":"bookmark code invalid"}`,
		},
		{
			name:             "returns internal server error when service fails",
			body:             requestBody,
			withClaims:       true,
			callService:      true,
			serviceErr:       errors.New("service error"),
			wantStatus:       http.StatusInternalServerError,
			wantResponseBody: `{"message":"Processing error"}`,
		},
		{
			name:             "returns bad request for invalid URL",
			body:             `{"description":"Go documentation","url":"not-a-url"}`,
			wantStatus:       http.StatusBadRequest,
			wantResponsePart: `"message":"Invalid input"`,
		},
		{
			name:             "returns unauthorized when claims are missing",
			body:             requestBody,
			wantStatus:       http.StatusUnauthorized,
			wantResponseBody: `{"error":"Invalid token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(
				http.MethodPost,
				"/v1/bookmarks",
				bytes.NewBufferString(tt.body),
			)
			ctx.Request.Header.Set("Content-Type", "application/json")
			if tt.withClaims {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
			}

			service := mocks.NewBookmarkService(t)
			if tt.callService {
				service.On("AddBookmark", ctx, "Go documentation", "https://go.dev/doc/", userID).
					Return(tt.serviceResult, tt.serviceErr).
					Once()
			}

			NewBookmarkHandler(service).AddBookmark(ctx)

			assert.Equal(t, tt.wantStatus, recorder.Code)
			if tt.wantResponsePart != "" {
				assert.Contains(t, recorder.Body.String(), tt.wantResponsePart)
				return
			}
			assert.JSONEq(t, tt.wantResponseBody, recorder.Body.String())
		})
	}
}
