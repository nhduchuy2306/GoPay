package auth

import (
	"gopay/internal/user"
	"gopay/internal/utils"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	Login(request user.LoginRequest) user.LoginResponse
	Register(request user.User) user.User
}

type service struct {
	userService user.Service
}

func NewService(userService user.Service) Service {
	return &service{userService: userService}
}

func (s *service) Login(request user.LoginRequest) user.LoginResponse {
	loginUser, ok := s.userService.GetByEmail(request.Email)
	if !ok {
		log.Fatal("Can not login")
	}
	hashPassword := loginUser.Password
	if ok := utils.MatchPassword(hashPassword, request.Password); !ok {
		log.Fatal("Password not correct")
	}
	token, err := s.generateToken(loginUser)
	if err != nil {
		log.Fatal(err)
	}
	return user.LoginResponse{
		Token: token,
	}
}

func (s *service) Register(request user.User) user.User {
	registerUser := s.userService.Create(request)
	return registerUser
}

func (s *service) VerifyToken(token string) (*jwt.Token, error) {
	secret := []byte(os.Getenv("SECRET_KEY"))
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func (s *service) generateToken(user user.User) (string, error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("SECRET_KEY"))
	return token.SignedString(secret)
}
