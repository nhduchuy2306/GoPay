package user

import (
	"gopay/internal/common"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	common.BaseModel
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     Role   `gorm:"not null" json:"role"`
}

type LoginRequest struct {
	Email    string
	Password string
}
