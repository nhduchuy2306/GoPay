package wallet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetAll() ([]Wallet, error)
	GetByID(id string) (Wallet, error)
	GetByUserID(userID string) (Wallet, error)
	Delete(id string) error
	Create(body Wallet) (Wallet, error)
	Update(id string, body Wallet) (Wallet, error)
	Deposit(id string, amount int64) (Wallet, error)
	Withdraw(id string, amount int64) (Wallet, error)
	Transfer(req TransferRequest) (TransferResult, error)
	History(id string) ([]WalletTransactionHistory, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type TransferResult struct {
	From Wallet `json:"from_wallet"`
	To   Wallet `json:"to_wallet"`
}

func (s *service) GetAll() ([]Wallet, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id string) (Wallet, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByUserID(userID string) (Wallet, error) {
	return s.repo.GetByUserID(userID)
}

func (s *service) Create(body Wallet) (Wallet, error) {
	if body.ID == "" {
		body.ID = uuid.NewString()
	}
	if body.Currency == "" {
		body.Currency = "VND"
	}
	if body.CreatedAt.IsZero() {
		body.CreatedAt = time.Now()
	}
	if body.Version == 0 {
		body.Version = 1
	}
	return s.repo.Create(body)
}

func (s *service) Update(id string, body Wallet) (Wallet, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return Wallet{}, err
	}
	if body.UserID != "" {
		current.UserID = body.UserID
	}
	if body.Currency != "" {
		current.Currency = body.Currency
	}
	if body.Balance != 0 {
		current.Balance = body.Balance
	}
	current.Version++
	return s.repo.Update(current)
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *service) Deposit(id string, amount int64) (Wallet, error) {
	if amount <= 0 {
		return Wallet{}, errors.New("amount must be greater than zero")
	}
	return s.repo.AdjustBalance(id, amount)
}

func (s *service) Withdraw(id string, amount int64) (Wallet, error) {
	if amount <= 0 {
		return Wallet{}, errors.New("amount must be greater than zero")
	}
	return s.repo.AdjustBalance(id, -amount)
}

func (s *service) Transfer(req TransferRequest) (TransferResult, error) {
	if req.Amount <= 0 {
		return TransferResult{}, errors.New("amount must be greater than zero")
	}
	from, err := s.repo.AdjustBalance(req.FromWalletID, -req.Amount)
	if err != nil {
		return TransferResult{}, err
	}
	to, err := s.repo.AdjustBalance(req.ToWalletID, req.Amount)
	if err != nil {
		_, _ = s.repo.AdjustBalance(req.FromWalletID, req.Amount)
		return TransferResult{}, err
	}
	return TransferResult{From: from, To: to}, nil
}

func (s *service) History(id string) ([]WalletTransactionHistory, error) {
	return s.repo.History(id)
}
