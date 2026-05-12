package wallet

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() []Wallet
	GetByID(id string) (Wallet, bool)
	Delete(id string) (Wallet, bool)
	Create(body Wallet) Wallet
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	if err := db.AutoMigrate(&Wallet{}); err != nil {
		log.Fatal(err)
	}
	return &repo{db: db}
}

func (r *repo) GetAll() []Wallet {
	var wallets []Wallet
	if err := r.db.Find(&wallets).Error; err != nil {
		return nil
	}
	return wallets
}

func (r *repo) GetByID(id string) (Wallet, bool) {
	var wallet Wallet
	if err := r.db.First(&wallet, id).Error; err != nil {
		return Wallet{}, false
	}
	return wallet, true
}

func (r *repo) Create(body Wallet) Wallet {
	if err := r.db.Create(&body).Error; err != nil {
		return Wallet{}
	}
	return body
}

func (r *repo) Delete(id string) (Wallet, bool) {
	var wallet Wallet
	if err := r.db.Delete(&wallet, id).Error; err != nil {
		return Wallet{}, false
	}
	return wallet, true
}
