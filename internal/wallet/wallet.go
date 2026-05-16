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

type WithdrawRequest struct {
	Amount int64 `json:"amount,omitempty"`
}

type TransferRequest struct {
	FromWalletID   string `json:"from_wallet_id"`
	ToWalletID     string `json:"to_wallet_id"`
	Amount         int64  `json:"amount"`
	Description    string `json:"description,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

type CreateRequest struct {
	UserID   string `json:"user_id"`
	Currency string `json:"currency"`
	Balance  int64  `json:"balance,omitempty"`
}
