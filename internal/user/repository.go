package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() []User
	GetByID(id string) (User, bool)
	GetByEmail(email string) (User, bool)
	GetByEmailAndPassword(email string, password string) (User, bool)
	Create(u User) User
	Update(id string, u User) (User, bool)
	Delete(id string) bool
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
	return &repo{db: db}
}

func (r *repo) GetAll() []User {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil
	}
	return users
}

func (r *repo) GetByID(id string) (User, bool) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return User{}, false
	}
	return user, true
}

func (r *repo) GetByEmail(email string) (User, bool) {
	var user User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return User{}, false
	}
	return user, true
}

func (r *repo) GetByEmailAndPassword(email string, password string) (User, bool) {
	var user User
	if err := r.db.First(&user, "email = ? and password = ?", email, password).Error; err != nil {
		return User{}, false
	}
	return user, true
}

func (r *repo) Create(u User) User {
	if err := r.db.Create(&u).Error; err != nil {
		return User{}
	}
	return u
}

func (r *repo) Update(id string, u User) (User, bool) {
	var existing User
	if err := r.db.First(&existing, "id = ?", id).Error; err != nil {
		return User{}, false
	}

	existing.Email = u.Email
	existing.Password = u.Password
	existing.Role = u.Role

	if err := r.db.Save(&existing).Error; err != nil {
		return User{}, false
	}
	return existing, true
}

func (r *repo) Delete(id string) bool {
	if err := r.db.Delete(&User{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}
