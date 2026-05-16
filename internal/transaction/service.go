package transaction

import (
	"gopay/internal/notification"
	"gopay/internal/wallet"
	"log"
)

type Service interface {
	GetAll() ([]Transaction, error)
	GetByID(id string) (Transaction, error)
	ListByWalletID(walletID string) ([]Transaction, error)
	Create(req CreateRequest) (Transaction, error)
	UpdateStatus(id string, status string) (Transaction, error)
}

type service struct {
	repo                Repository
	walletService       wallet.Service
	notificationService *notification.Service
}

func NewService(repo Repository, walletService wallet.Service, notificationService *notification.Service) Service {
	return &service{repo: repo, walletService: walletService, notificationService: notificationService}
}

func (s *service) GetAll() ([]Transaction, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id string) (Transaction, error) {
	return s.repo.GetByID(id)
}

func (s *service) ListByWalletID(walletID string) ([]Transaction, error) {
	return s.repo.ListByWalletID(walletID)
}

func (s *service) Create(req CreateRequest) (Transaction, error) {
	tx, err := s.repo.Create(req)
	if err != nil {
		return Transaction{}, err
	}

	if s.notificationService != nil && s.walletService != nil {
		go func() {
			toWallet, walletErr := s.walletService.GetByID(req.ToWalletID)
			if walletErr != nil {
				log.Println("transfer email wallet lookup failed:", walletErr)
				return
			}
			if toWallet.User.Email == "" {
				return
			}
			if err := s.notificationService.SendTransferSuccessEmail(toWallet.User.Email, notification.TransferSuccessEmailData{
				Amount:        tx.Amount,
				Currency:      toWallet.Currency,
				TransactionID: tx.ID,
				FromWalletID:  tx.FromWalletID,
				ToWalletID:    tx.ToWalletID,
				Description:   tx.Description,
			}); err != nil {
				log.Println("transfer success email send failed:", err)
			}
		}()
	}

	return tx, nil
}

func (s *service) UpdateStatus(id string, status string) (Transaction, error) {
	return s.repo.UpdateStatus(id, status)
}
