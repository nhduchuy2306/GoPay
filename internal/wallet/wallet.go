package wallet

import (
	"gopay/internal/common"
	"gopay/internal/user"
	"time"
)

type Wallet struct {
	common.BaseModel
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	User      user.User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Balance   int64     `gorm:"not null" json:"balance"`
	Currency  string    `gorm:"not null" json:"currency"`
	Version   int       `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DepositRequest struct {
	Amount int64 `json:"amount,omitempty"`
}
