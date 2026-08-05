package bookmark

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/service/bookmark/mocks"
	"github.com/stretchr/testify/require"
)

const (
	updateBookmarkID     = "bookmark-123"
	updateBookmarkUserID = "user-123"
	updateDescription    = "Updated description"
	updateURL            = "https://go.dev/doc/"
)

type updateBookmarkHandlerTestCase struct {
	name             string
	body             string
	query            string
	configureService bool
	serviceErr       error
	wantStatus       int
	wantResponse     string
	wantResponsePart string
}

func TestBookmarkHandler_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []updateBookmarkHandlerTestCase{
		{
			name:             "updates bookmark",
			body:             `{"description":"Updated description","url":"https://go.dev/doc/"}`,
			query:            "?id=" + updateBookmarkID,
			configureService: true,
			wantStatus:       http.StatusOK,
			wantResponse:     `{"message":"Success"}`,
		},
		{
			name:             "returns internal server error when service fails",
			body:             `{"description":"Updated description","url":"https://go.dev/doc/"}`,
			query:            "?id=" + updateBookmarkID,
			configureService: true,
			serviceErr:       errors.New("update bookmark failed"),
			wantStatus:       http.StatusInternalServerError,
			wantResponse:     `{"message":"Fail to update bookmark"}`,
		},
		{
			name:             "returns bad request for invalid URL",
			body:             `{"description":"Updated description","url":"not-a-url"}`,
			query:            "?id=" + updateBookmarkID,
			wantStatus:       http.StatusBadRequest,
			wantResponsePart: `"message":"Invalid input"`,
		},
		{
			name:             "returns bad request when bookmark ID is missing",
			body:             `{"description":"Updated description","url":"https://go.dev/doc/"}`,
			wantStatus:       http.StatusBadRequest,
			wantResponsePart: `"message":"Invalid input"`,
		},
		{
			name:             "returns bad request for malformed JSON",
			body:             `{"description":`,
			query:            "?id=" + updateBookmarkID,
			wantStatus:       http.StatusBadRequest,
			wantResponsePart: `"message":"Invalid input"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx, recorder := newUpdateBookmarkContext(testCase.body, testCase.query)
			service := mocks.NewBookmarkService(t)
			if testCase.configureService {
				service.On("UpdateBookmark", ctx, updateBookmarkUserID, updateBookmarkID, updateDescription, updateURL).
					Return(testCase.serviceErr).Once()
			}

			NewBookmarkHandler(service).UpdateBookmark(ctx)

			assertUpdateBookmarkResponse(t, recorder, testCase.wantStatus, testCase.wantResponse, testCase.wantResponsePart)
		})
	}
}

func newUpdateBookmarkContext(body, query string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(
		http.MethodPut,
		"/v1/bookmarks"+query,
		bytes.NewBufferString(body),
	)
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("claims", jwt.MapClaims{"sub": updateBookmarkUserID})
	return ctx, recorder
}

func assertUpdateBookmarkResponse(t *testing.T, recorder *httptest.ResponseRecorder, status int, body, bodyPart string) {
	t.Helper()
	require.Equal(t, status, recorder.Code)
	if bodyPart != "" {
		require.Contains(t, recorder.Body.String(), bodyPart)
		return
	}
	require.JSONEq(t, body, recorder.Body.String())
}
