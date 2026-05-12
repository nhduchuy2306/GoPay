package middleware

import (
	"gopay/internal/auth/token"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	IgnoreRequestUri = []string{"/login", "/register"}
)

func BearerMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestUri := ctx.Request.RequestURI
		isIgnore := false
		for _, path := range IgnoreRequestUri {
			if strings.HasSuffix(requestUri, path) {
				isIgnore = true
				break
			}
		}
		if isIgnore {
			ctx.Next()
			return
		}
		bearerToken, message, ok := token.GetBearerToken(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(401, gin.H{"error": message})
			return
		}
		ctx.Set("token", bearerToken)
		ctx.Next()
	}
}
