package webhook

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	ListEndpoints() ([]Endpoint, error)
	ListMyEndpoints(userID string) ([]Endpoint, error)
	GetEndpoint(id string) (Endpoint, error)
	CreateEndpoint(userID string, req CreateEndpointRequest) (Endpoint, error)
	UpdateEndpoint(id string, req UpdateEndpointRequest) (Endpoint, error)
	DeleteEndpoint(id string) error
	ListDeliveries(endpointID string) ([]Delivery, error)
	GetDelivery(id string) (Delivery, error)
}

type service struct {
	repo   Repository
	worker Worker
}

func NewService(repo Repository, worker Worker) Service {
	return &service{repo: repo, worker: worker}
}

func (s *service) ListEndpoints() ([]Endpoint, error) {
	return s.repo.GetAll()
}

func (s *service) ListMyEndpoints(userID string) ([]Endpoint, error) {
	return s.repo.GetByUserID(userID)
}

func (s *service) GetEndpoint(id string) (Endpoint, error) {
	return s.repo.GetByID(id)
}

func (s *service) CreateEndpoint(userID string, req CreateEndpointRequest) (Endpoint, error) {
	if req.URL == "" {
		return Endpoint{}, errors.New("url is required")
	}
	endpoint := Endpoint{
		UserID:   userID,
		URL:      req.URL,
		Secret:   req.Secret,
		IsActive: req.IsActive,
	}
	if endpoint.Secret == "" {
		endpoint.Secret = uuid.NewString()
	}
	if endpoint.CreatedAt.IsZero() {
		endpoint.CreatedAt = time.Now()
	}
	return s.repo.Create(endpoint)
}

func (s *service) UpdateEndpoint(id string, req UpdateEndpointRequest) (Endpoint, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return Endpoint{}, err
	}
	if req.URL != "" {
		current.URL = req.URL
	}
	if req.Secret != "" {
		current.Secret = req.Secret
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}
	return s.repo.Update(current)
}

func (s *service) DeleteEndpoint(id string) error {
	return s.repo.Delete(id)
}

func (s *service) ListDeliveries(endpointID string) ([]Delivery, error) {
	return s.repo.GetDeliveriesByEndpointID(endpointID)
}

func (s *service) GetDelivery(id string) (Delivery, error) {
	return s.repo.GetDeliveryByID(id)
}
