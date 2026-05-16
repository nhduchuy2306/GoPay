package user

import (
	"gopay/internal/common"
	"gopay/internal/utils"
	"gopay/internal/utils/enum"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetUsers() ([]User, error)
	GetByID(id string) (User, error)
	GetByEmail(email string) (User, error)
	GetByEmailAndPassword(email string, password string) (User, error)
	Create(u User) (User, error)
	Update(id string, u User) (User, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) GetUsers() ([]User, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id string) (User, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByEmail(email string) (User, error) {
	return s.repo.GetByEmail(email)
}

func (s *service) GetByEmailAndPassword(email string, password string) (User, error) {
	return s.repo.GetByEmailAndPassword(email, password)
}

func (s *service) Create(u User) (User, error) {
	base := common.BaseModel{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}
	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return User{}, err
	}
	user := User{
		BaseModel: base,
		Email:     u.Email,
		Password:  hashPassword,
		Role:      enum.RoleUser,
	}
	return s.repo.Create(user)
}

func (s *service) Update(id string, u User) (User, error) {
	if u.Password != "" {
		hashPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return User{}, err
		}
		u.Password = hashPassword
	}
	return s.repo.Update(id, u)
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}
