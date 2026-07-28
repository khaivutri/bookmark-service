package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/pkg/jwtutils"
)

type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtVal jwtutils.JWTValidator
}

func NewJWTAuth(jwtVal jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{jwtVal: jwtVal}
}

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
