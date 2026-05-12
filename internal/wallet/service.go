package wallet

type Service interface {
	GetAll() []Wallet
	GetByID(id string) (Wallet, bool)
	Delete(id string) (Wallet, bool)
	Create(body Wallet) Wallet
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() []Wallet {
	return s.repo.GetAll()
}

func (s *service) GetByID(id string) (Wallet, bool) {
	return s.repo.GetByID(id)
}

func (s *service) Create(body Wallet) Wallet {
	return s.repo.Create(body)
}

func (s *service) Delete(id string) (Wallet, bool) {
	return s.repo.Delete(id)
}
