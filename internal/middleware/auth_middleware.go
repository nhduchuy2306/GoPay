package middleware

import (
	"gopay/internal/auth/token"
	"gopay/internal/utils/enum"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken, ok := ctx.Get("token")
		if !ok {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Missing Bearer Token"})
			return
		}

		tokenStr, ok := bearerToken.(string)
		if !ok {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "Can not parse token"})
			return
		}

		parseToken, err := token.ParseToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(403, gin.H{"error": err})
			return
		}
		role := parseToken["role"]
		if role != enum.RoleAdmin {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "Can not access this resource"})
			return
		}
		ctx.Next()
	}
}
