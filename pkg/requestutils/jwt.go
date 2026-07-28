package requestutils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ( 
	ErrNoClaims = errors.New("token format invalid")
)
func GetUserIDFromRequest(ctx *gin	.Context) (string, error){
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"Invalid token"})
		ctx.Abort()
		return "", ErrNoClaims
	}

	tokenInfo, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"Invalid token"})
		ctx.Abort()
		return "", ErrNoClaims
	}

	uid, ok := tokenInfo["sub"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"Invalid token"})
		ctx.Abort()
		return "", ErrNoClaims
	}

	return uid, nil
}