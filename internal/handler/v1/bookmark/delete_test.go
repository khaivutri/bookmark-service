package bookmark

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/service/bookmark/mocks"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/require"
)

const deleteBookmarkID = "bookmark-123"

type deleteBookmarkHandlerTestCase struct {
	name             string
	query            string
	configureService bool
	serviceErr       error
	wantStatus       int
	wantResponse     string
	wantResponsePart string
}

func TestBookmarkHandler_DeleteBookmark(t *testing.T) {
	t.Parallel()

	testCases := []deleteBookmarkHandlerTestCase{
		{
			name:             "deletes bookmark",
			query:            "?id=" + deleteBookmarkID,
			configureService: true,
			wantStatus:       http.StatusOK,
			wantResponse:     `{"message":"Success"}`,
		},
		{
			name:             "returns not found when service cannot find bookmark",
			query:            "?id=" + deleteBookmarkID,
			configureService: true,
			serviceErr:       dbutils.ErrRecordNotFound,
			wantStatus:       http.StatusNotFound,
			wantResponse:     `{"message":"Bookmark not found"}`,
		},
		{
			name:             "returns internal server error when service fails",
			query:            "?id=" + deleteBookmarkID,
			configureService: true,
			serviceErr:       errors.New("delete bookmark failed"),
			wantStatus:       http.StatusInternalServerError,
			wantResponse:     `{"message":"Fail to delete bookmark"}`,
		},
		{
			name:             "returns bad request when bookmark ID is missing",
			wantStatus:       http.StatusBadRequest,
			wantResponsePart: `"message":"Invalid input"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, recorder := newDeleteBookmarkContext(tc.query)
			service := mocks.NewBookmarkService(t)
			if tc.configureService {
				service.On("DeleteBookmark", ctx, deleteBookmarkID).
					Return(tc.serviceErr).
					Once()
			}

			(&bookmarkHandler{svc: service}).DeleteBookmark(ctx)

			assertDeleteBookmarkResponse(t, recorder, tc.wantStatus, tc.wantResponse, tc.wantResponsePart)
		})
	}
}

func newDeleteBookmarkContext(query string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/v1/bookmarks"+query, nil)
	return ctx, recorder
}

func assertDeleteBookmarkResponse(t *testing.T, recorder *httptest.ResponseRecorder, status int, response, responsePart string) {
	t.Helper()
	require.Equal(t, status, recorder.Code)
	if responsePart != "" {
		require.Contains(t, recorder.Body.String(), responsePart)
		return
	}
	require.JSONEq(t, response, recorder.Body.String())
}
