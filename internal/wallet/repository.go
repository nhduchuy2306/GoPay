package wallet

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetAll() ([]Wallet, error)
	GetByID(id string) (Wallet, error)
	GetByUserID(userID string) (Wallet, error)
	Delete(id string) error
	Create(body Wallet) (Wallet, error)
	Update(body Wallet) (Wallet, error)
	AdjustBalance(id string, delta int64) (Wallet, error)
	History(id string) ([]WalletTransactionHistory, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	_ = db.AutoMigrate(&Wallet{})
	return &repo{db: db}
}

func (r *repo) GetAll() ([]Wallet, error) {
	var wallets []Wallet
	if err := r.db.Preload("User").Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *repo) GetByID(id string) (Wallet, error) {
	var wallet Wallet
	if err := r.db.Preload("User").First(&wallet, "id = ?", id).Error; err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

func (r *repo) GetByUserID(userID string) (Wallet, error) {
	var wallet Wallet
	if err := r.db.Preload("User").First(&wallet, "user_id = ?", userID).Error; err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

func (r *repo) Create(body Wallet) (Wallet, error) {
	if err := r.db.Create(&body).Error; err != nil {
		return Wallet{}, err
	}
	return body, nil
}

func (r *repo) Update(body Wallet) (Wallet, error) {
	if err := r.db.Save(&body).Error; err != nil {
		return Wallet{}, err
	}
	return body, nil
}

func (r *repo) Delete(id string) error {
	if err := r.db.Delete(&Wallet{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *repo) AdjustBalance(id string, delta int64) (Wallet, error) {
	var updated Wallet
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var wallet Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, "id = ?", id).Error; err != nil {
			return err
		}
		if wallet.Balance+delta < 0 {
			return errors.New("insufficient balance")
		}
		wallet.Balance += delta
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}
		updated = wallet
		return nil
	})
	if err != nil {
		return Wallet{}, err
	}
	return updated, nil
}

func (r *repo) History(id string) ([]WalletTransactionHistory, error) {
	var history []WalletTransactionHistory
	if err := r.db.Table("transactions").
		Select("id, from_wallet_id, to_wallet_id, amount, status, description, created_at").
		Where("from_wallet_id = ? OR to_wallet_id = ?", id, id).
		Order("created_at desc").
		Scan(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

type WalletTransactionHistory struct {
	ID           string    `json:"id"`
	FromWalletID string    `json:"from_wallet_id"`
	ToWalletID   string    `json:"to_wallet_id"`
	Amount       int64     `json:"amount"`
	Status       string    `json:"status"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
