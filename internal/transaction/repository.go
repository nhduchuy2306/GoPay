package transaction

import (
	"encoding/json"
	"errors"
	"gopay/internal/common"
	"gopay/internal/outbox"
	"gopay/internal/utils/enum"
	"gopay/internal/wallet"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Transaction, error)
	GetByID(id string) (Transaction, error)
	GetByIdempotencyKey(key string) (Transaction, error)
	ListByWalletID(walletID string) ([]Transaction, error)
	Create(req CreateRequest) (Transaction, error)
	UpdateStatus(id string, status string) (Transaction, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	_ = db.AutoMigrate(&Transaction{})
	return &repo{db: db}
}

func (r *repo) GetAll() ([]Transaction, error) {
	var txs []Transaction
	if err := r.db.Order("created_at desc").Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *repo) GetByID(id string) (Transaction, error) {
	var tx Transaction
	if err := r.db.First(&tx, "id = ?", id).Error; err != nil {
		return Transaction{}, err
	}
	return tx, nil
}

func (r *repo) GetByIdempotencyKey(key string) (Transaction, error) {
	var tx Transaction
	if err := r.db.First(&tx, "idempotency_key = ?", key).Error; err != nil {
		return Transaction{}, err
	}
	return tx, nil
}

func (r *repo) ListByWalletID(walletID string) ([]Transaction, error) {
	var txs []Transaction
	if err := r.db.
		Where("from_wallet_id = ? OR to_wallet_id = ?", walletID, walletID).
		Order("created_at desc").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *repo) Create(req CreateRequest) (Transaction, error) {
	if req.Amount <= 0 {
		return Transaction{}, errors.New("amount must be greater than zero")
	}
	if req.FromWalletID == "" || req.ToWalletID == "" {
		return Transaction{}, errors.New("wallet ids are required")
	}
	if req.FromWalletID == req.ToWalletID {
		return Transaction{}, errors.New("source and destination wallets must differ")
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = uuid.NewString()
	}

	var persisted Transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if existing, err := func() (Transaction, error) {
			var existing Transaction
			if err := tx.First(&existing, "idempotency_key = ?", req.IdempotencyKey).Error; err != nil {
				return Transaction{}, err
			}
			return existing, nil
		}(); err == nil {
			persisted = existing
			return nil
		}

		var fromWallet wallet.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fromWallet, "id = ?", req.FromWalletID).Error; err != nil {
			return err
		}
		var toWallet wallet.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&toWallet, "id = ?", req.ToWalletID).Error; err != nil {
			return err
		}
		if fromWallet.Balance < req.Amount {
			return errors.New("insufficient balance")
		}

		fromWallet.Balance -= req.Amount
		toWallet.Balance += req.Amount

		if err := tx.Save(&fromWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&toWallet).Error; err != nil {
			return err
		}

		persisted = Transaction{
			BaseModel:      common.BaseModel{ID: uuid.NewString()},
			IdempotencyKey: req.IdempotencyKey,
			FromWalletID:   req.FromWalletID,
			ToWalletID:     req.ToWalletID,
			Amount:         req.Amount,
			Status:         enum.Success,
			Description:    req.Description,
		}
		if err := tx.Create(&persisted).Error; err != nil {
			return err
		}

		payload, err := json.Marshal(map[string]any{
			"transaction_id": persisted.ID,
			"amount":         persisted.Amount,
			"from_wallet_id": persisted.FromWalletID,
			"to_wallet_id":   persisted.ToWalletID,
			"status":         persisted.Status,
			"description":    persisted.Description,
		})
		if err != nil {
			return err
		}

		outboxEvent := outbox.Outbox{
			BaseModel:   common.BaseModel{ID: uuid.NewString()},
			AggregateID: persisted.ID,
			EventType:   enum.TransactionCompleted,
			Payload:     datatypes.JSON(payload),
		}
		return tx.Create(&outboxEvent).Error
	})
	if err != nil {
		return Transaction{}, err
	}
	return persisted, nil
}

func (r *repo) UpdateStatus(id string, status string) (Transaction, error) {
	var txRecord Transaction
	if err := r.db.First(&txRecord, "id = ?", id).Error; err != nil {
		return Transaction{}, err
	}
	txRecord.Status = status
	if err := r.db.Save(&txRecord).Error; err != nil {
		return Transaction{}, err
	}
	return txRecord, nil
}
