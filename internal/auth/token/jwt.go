package token

import (
	"errors"
	"gopay/internal/utils/constants"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenUser interface {
	GetID() string
	GetEmail() string
	GetRole() string
}

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

func VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	secret := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func GenerateToken(user TokenUser) (string, error) {
	claims := jwt.MapClaims{
		"userId": user.GetID(),
		"email":  user.GetEmail(),
		"role":   user.GetRole(),
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	}
	newTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("SECRET_KEY"))
	return newTokenClaims.SignedString(secret)
}
