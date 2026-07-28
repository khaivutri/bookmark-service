package requestutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(c *gin.Context)
		expectedUserID string
		expectedErr    error
		expectedStatus int
	}{
		{
			name: "Success - Valid claims with sub field",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{
					"sub": "user_123456",
				})
			},
			expectedUserID: "user_123456",
			expectedErr:    nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failure - No claims in context",
			setupContext: func(c *gin.Context) {
			},
			expectedUserID: "",
			expectedErr:    ErrNoClaims,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Failure - Claims is not jwt.MapClaims type",
			setupContext: func(c *gin.Context) {
				c.Set("claims", "invalid_claims_type_string")
			},
			expectedUserID: "",
			expectedErr:    ErrNoClaims,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Failure - Sub field missing in claims",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{
					"email": "test@example.com",
				})
			},
			expectedUserID: "",
			expectedErr:    ErrNoClaims,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Failure - Sub field is not a string",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{
					"sub": 12345, // int instead of string
				})
			},
			expectedUserID: "",
			expectedErr:    ErrNoClaims,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setupContext(c)

			userID, err := GetUserIDFromRequest(c)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUserID, userID)
			
			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.True(t, c.IsAborted())
			}
		})
	}
}