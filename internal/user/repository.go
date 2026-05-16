package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]User, error)
	GetByID(id string) (User, error)
	GetByEmail(email string) (User, error)
	GetByEmailAndPassword(email string, password string) (User, error)
	Create(u User) (User, error)
	Update(id string, u User) (User, error)
	Delete(id string) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	_ = db.AutoMigrate(&User{})
	return &repo{db: db}
}

func (r *repo) GetAll() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repo) GetByID(id string) (User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repo) GetByEmail(email string) (User, error) {
	var user User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repo) GetByEmailAndPassword(email string, password string) (User, error) {
	var user User
	if err := r.db.First(&user, "email = ? and password = ?", email, password).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repo) Create(u User) (User, error) {
	if err := r.db.Create(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *repo) Update(id string, u User) (User, error) {
	var existing User
	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return User{}, err
	}

	if u.Email != "" {
		existing.Email = u.Email
	}
	if u.Password != "" {
		existing.Password = u.Password
	}
	if u.Role != "" {
		existing.Role = u.Role
	}

	if err := r.db.Save(&existing).Error; err != nil {
		return User{}, err
	}
	return existing, nil
}

func (r *repo) Delete(id string) error {
	if err := r.db.Delete(&User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

var ErrUserNotFound = errors.New("user not found")
