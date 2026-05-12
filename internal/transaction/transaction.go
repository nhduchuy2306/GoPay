package transaction

import (
	"gopay/internal/common"
	"gopay/internal/wallet"
	"time"
)

type TransactionStatus string

const (
	StatusPending TransactionStatus = "pending"
	StatusSuccess TransactionStatus = "success"
	StatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	common.BaseModel
	IdempotencyKey string            `gorm:"not null;uniqueIndex" json:"idempotency_key"`
	FromWalletID   string            `gorm:"type:uuid;not null" json:"from_wallet_id"`
	ToWalletID     string            `gorm:"type:uuid;not null" json:"to_wallet_id"`
	FromWallet     wallet.Wallet     `gorm:"foreignKey:FromWalletID;references:ID" json:"from_wallet,omitempty"`
	ToWallet       wallet.Wallet     `gorm:"foreignKey:ToWalletID;references:ID" json:"to_wallet,omitempty"`
	Amount         int64             `gorm:"not null" json:"amount"`
	Status         TransactionStatus `gorm:"not null" json:"status"`
	Description    string            `json:"description,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
}
