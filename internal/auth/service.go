package auth

import (
	"gopay/internal/user"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest = user.LoginRequest
type UserService = user.Service
type User = user.User

type Service interface {
	Login(request LoginRequest) User
	Register(request User) User
}

type service struct {
	userService UserService
}

func NewService(userService UserService) Service {
	return &service{userService: userService}
}

func (s *service) Login(request LoginRequest) User {
	loginUser, ok := s.userService.GetByEmailAndPassword(request.Email, request.Password)
	if !ok {
		panic("Can not login")
	}
	return loginUser
}

func (s *service) Register(request User) User {
	registerUser := s.userService.Create(request)
	return registerUser
}

func (s *service) VerifyToken(token string) (*jwt.Token, error) {
	secret := []byte(os.Getenv("SECRET_KEY"))
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func (s *service) generateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("SECRET_KEY"))
	return token.SignedString(secret)
}
