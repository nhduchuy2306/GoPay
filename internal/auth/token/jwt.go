package token

import (
	"errors"
	"gopay/internal/utils/constants"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetBearerToken(ctx *gin.Context) (string, string, bool) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", "Authorization header missing", false
	}

	if !strings.HasPrefix(authHeader, constants.BearerPrefix) {
		return "", "Invalid authorization format", false
	}

	token := strings.TrimPrefix(authHeader, constants.BearerPrefix)

	return token, "Invalid authorization format", true
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("secret key not set")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
