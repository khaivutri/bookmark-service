package bookmark

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/model"
	servicebookmark "github.com/khaivutri/bookmark-service/internal/service/bookmark"
	"github.com/khaivutri/bookmark-service/internal/service/bookmark/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookmarkHandler_GetBookmarks(t *testing.T) {
	t.Parallel()

	const (
		userID = "user-123"
		page   = 2
		limit  = 10
	)

	bookmarks := []*model.Bookmark{{
		Description: "Go documentation",
		URL:         "https://go.dev/doc/",
		Code:        "go-docs",
		UserID:      userID,
	}}

	tests := []struct {
		name             string
		query            string
		withClaims       bool
		callService      bool
		page             int
		limit            int
		serviceResult    *servicebookmark.GetBookmarkResponse
		serviceErr       error
		wantStatus       int
		wantResponseBody string
		wantResponsePart string
	}{
		{
			name:        "returns paginated bookmarks",
			query:       "?page=2&limit=10",
			withClaims:  true,
			callService: true,
			page:        page,
			limit:       limit,
			serviceResult: &servicebookmark.GetBookmarkResponse{
				Bookmarks: bookmarks,
				Total:     21,
			},
			wantStatus:       http.StatusOK,
			wantResponseBody: `{"data":[{"id":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","description":"Go documentation","url":"https://go.dev/doc/","code":"go-docs","user_id":"user-123"}],"pagination":{"page":2,"limit":10,"total":21}}`,
		},
		{
			name:             "returns bad request for invalid pagination",
			query:            "?page=0&limit=10",
			wantStatus:       http.StatusBadRequest,
			wantResponsePart: `"message":"Invalid input"`,
		},
		{
			name:             "returns unauthorized without claims",
			query:            "?page=1&limit=10",
			wantStatus:       http.StatusUnauthorized,
			wantResponseBody: `{"error":"Invalid token"}`,
		},
		{
			name:             "returns internal server error when service fails",
			query:            "?page=1&limit=10",
			withClaims:       true,
			callService:      true,
			page:             1,
			limit:            10,
			serviceErr:       errors.New("get bookmarks failed"),
			wantStatus:       http.StatusInternalServerError,
			wantResponseBody: `{"message":"Processing error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, recorder := newGetBookmarksContext(test.query, test.withClaims, userID)
			service := mocks.NewBookmarkService(t)
			if test.callService {
				service.On("GetBookmarks", ctx, userID, test.page, test.limit).
					Return(test.serviceResult, test.serviceErr).Once()
			}

			NewBookmarkHandler(service).GetBookmarks(ctx)

			assertGetBookmarksResponse(t, recorder, test.wantStatus, test.wantResponseBody, test.wantResponsePart)
		})
	}
}

func newGetBookmarksContext(query string, withClaims bool, userID string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks"+query, nil)
	if withClaims {
		ctx.Set("claims", jwt.MapClaims{"sub": userID})
	}
	return ctx, recorder
}

func assertGetBookmarksResponse(t *testing.T, recorder *httptest.ResponseRecorder, status int, body, part string) {
	t.Helper()
	assert.Equal(t, status, recorder.Code)
	if part != "" {
		assert.Contains(t, recorder.Body.String(), part)
		return
	}
	require.JSONEq(t, body, recorder.Body.String())
}
