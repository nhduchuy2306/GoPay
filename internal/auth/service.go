package auth

import (
	"errors"
	"gopay/internal/auth/token"
	"gopay/internal/notification"
	"gopay/internal/user"
	"gopay/internal/utils"
	"gopay/internal/wallet"
	"log"
)

type Service interface {
	Login(request user.LoginRequest) (user.LoginResponse, error)
	Register(request user.RegisterRequest) (user.User, error)
	Me(token string) (user.User, error)
	Refresh(token string) (string, error)
	Logout(token string) error
}

type service struct {
	userService         user.Service
	walletService       wallet.Service
	notificationService *notification.Service
}

func NewService(userService user.Service, walletService wallet.Service, notificationService *notification.Service) Service {
	return &service{
		userService:         userService,
		walletService:       walletService,
		notificationService: notificationService,
	}
}

func (s *service) Login(request user.LoginRequest) (user.LoginResponse, error) {
	loginUser, err := s.userService.GetByEmail(request.Email)
	if err != nil {
		return user.LoginResponse{}, err
	}
	hashPassword := loginUser.Password
	if ok := utils.MatchPassword(hashPassword, request.Password); !ok {
		return user.LoginResponse{}, errors.New("password not correct")
	}
	generateToken, err := token.GenerateToken(&loginUser)
	if err != nil {
		return user.LoginResponse{}, err
	}
	return user.LoginResponse{
		Token: generateToken,
	}, nil
}

func (s *service) Register(request user.RegisterRequest) (user.User, error) {
	registerUser, err := s.userService.Create(user.User{Email: request.Email, Password: request.Password})
	if err != nil {
		return user.User{}, err
	}
	if s.walletService != nil {
		_, _ = s.walletService.Create(wallet.Wallet{UserID: registerUser.ID, Balance: 0, Currency: "VND"})
	}
	if s.notificationService != nil {
		go func(email string) {
			if err := s.notificationService.SendWelcomeEmail(email); err != nil {
				log.Println("welcome email send failed:", err)
			}
		}(registerUser.Email)
	}
	return registerUser, nil
}

func (s *service) Me(tokenStr string) (user.User, error) {
	claims, err := token.VerifyToken(tokenStr)
	if err != nil {
		return user.User{}, err
	}
	userID, _ := claims["user_id"].(string)
	if userID == "" {
		email, _ := claims["email"].(string)
		if email == "" {
			return user.User{}, errors.New("invalid token claims")
		}
		return s.userService.GetByEmail(email)
	}
	return s.userService.GetByID(userID)
}

func (s *service) Refresh(tokenStr string) (string, error) {
	claims, err := token.VerifyToken(tokenStr)
	if err != nil {
		return "", err
	}
	usr, err := s.Me(tokenStr)
	if err != nil {
		return "", err
	}
	_ = claims
	return token.GenerateToken(&usr)
}

func (s *service) Logout(tokenStr string) error {
	_, err := token.VerifyToken(tokenStr)
	return err
}
