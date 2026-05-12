package user

import (
	"gopay/internal/common"
	"gopay/internal/utils"
)

type Service interface {
	GetUsers() []User
	GetByID(id string) (User, bool)
	GetByEmail(email string) (User, bool)
	GetByEmailAndPassword(email string, password string) (User, bool)
	Create(u User) User
	Update(id string, u User) (User, bool)
	Delete(id string) bool
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s service) GetUsers() []User {
	return s.repo.GetAll()
}

func (s service) GetByID(id string) (User, bool) {
	return s.repo.GetByID(id)
}

func (s service) GetByEmail(email string) (User, bool) {
	return s.repo.GetByEmail(email)
}

func (s service) GetByEmailAndPassword(email string, password string) (User, bool) {
	return s.repo.GetByEmailAndPassword(email, password)
}

func (s service) Create(u User) User {
	base := common.BaseModel{}
	base.BeforeCreate()
	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		panic(err)
	}
	user := User{
		BaseModel: base,
		Email:     u.Email,
		Password:  hashPassword,
		Role:      RoleUser,
	}
	return s.repo.Create(user)
}

func (s service) Update(id string, u User) (User, bool) {
	return s.repo.Update(id, u)
}

func (s service) Delete(id string) bool {
	return s.repo.Delete(id)
}
