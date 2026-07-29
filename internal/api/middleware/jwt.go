package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/pkg/jwtutils"
)

// JWTAuther defines the interface for creating a JWT authentication middleware.
type JWTAuther interface {
	// JWTAuth returns a Gin handler function that enforces JWT authentication.
	JWTAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtVal jwtutils.JWTValidator
}

// NewJWTAuth initializes and returns a JWTAuther instance using a JWTValidator.
func NewJWTAuth(jwtVal jwtutils.JWTValidator) JWTAuther {
	return &jwtAuth{jwtVal: jwtVal}
}

// JWTAuth returns a Gin handler function that enforces JWT authentication.
func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from header
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Validate token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		
		tokenString := parts[1]

		tokenClaims, err := j.jwtVal.ValidateJWT(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return 
		}

		// Set Claims to Context -> next handler can read 
		ctx.Set("claims", tokenClaims)


		ctx.Next()
	}
}
