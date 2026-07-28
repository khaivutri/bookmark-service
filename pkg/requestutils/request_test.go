package requestutils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bindInputTestRequest struct {
	Body      string `json:"body"`
	Path      string `uri:"path"`
	Query     int    `form:"query"`
	Header    int    `header:"X-Test-Header"`
	Validated string `json:"validated" form:"validated" validate:"required"`
}

func newBindTestContext(t *testing.T, method, target, body string, headers map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, target, strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		ctx.Request.Header.Set(key, value)
	}
	ctx.Params = gin.Params{{Key: "path", Value: "route-path"}}

	return ctx, recorder
}

func TestBindInputFromResquest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		method        string
		target        string
		body          string
		headers       map[string]string
		want           *bindInputTestRequest
		wantErr       bool
		wantStatus    int
		wantAborted   bool
	}{
		{
			name:    "binds JSON and overrides values from URI query and header",
			method:  http.MethodPost,
			target:  "/links/route-path?query=42",
			body:    `{"body":"json-body","path":"json-path","query":1,"header":2,"validated":"json-valid"}`,
			headers: map[string]string{"X-Test-Header": "7"},
			want: &bindInputTestRequest{
				Body: "json-body", Path: "route-path", Query: 42, Header: 7, Validated: "json-valid",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "GET skips malformed JSON and binds URI query and header",
			method:  http.MethodGet,
			target:  "/links/route-path?query=42&validated=query-valid",
			body:    `{"invalid-json"`,
			headers: map[string]string{"X-Test-Header": "7"},
			want: &bindInputTestRequest{
				Path: "route-path", Query: 42, Header: 7, Validated: "query-valid",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "ignores an empty POST body",
			method:  http.MethodPost,
			target:  "/links/route-path?query=42&validated=query-valid",
			headers: map[string]string{"X-Test-Header": "7"},
			want: &bindInputTestRequest{
				Path: "route-path", Query: 42, Header: 7, Validated: "query-valid",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "returns bad request for malformed JSON",
			method:     http.MethodPost,
			target:     "/links/route-path?query=42",
			body:       `{"body":`,
			headers:    map[string]string{"X-Test-Header": "7"},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
			wantAborted: true,
		},
		{
			name:       "returns bad request for an invalid query value",
			method:     http.MethodGet,
			target:     "/links/route-path?query=not-a-number&validated=query-valid",
			headers:    map[string]string{"X-Test-Header": "7"},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
			wantAborted: true,
		},
		{
			name:       "returns bad request for an invalid header value",
			method:     http.MethodGet,
			target:     "/links/route-path?query=42&validated=query-valid",
			headers:    map[string]string{"X-Test-Header": "not-a-number"},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
			wantAborted: true,
		},
		{
			name:       "returns bad request when final validation fails",
			method:     http.MethodGet,
			target:     "/links/route-path?query=42",
			headers:    map[string]string{"X-Test-Header": "7"},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
			wantAborted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, recorder := newBindTestContext(t, tt.method, tt.target, tt.body, tt.headers)

			got, err := BindInputFromResquest[bindInputTestRequest](ctx)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
				assert.Equal(t, tt.wantStatus, recorder.Code)
				assert.Equal(t, tt.wantAborted, ctx.IsAborted())
				assert.Contains(t, recorder.Body.String(), `"message":"Invalid input"`)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantStatus, recorder.Code)
			assert.False(t, ctx.IsAborted())
		})
	}
}
