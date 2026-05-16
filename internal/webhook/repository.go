package webhook

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Endpoint, error)
	GetActive() ([]Endpoint, error)
	GetByUserID(userID string) ([]Endpoint, error)
	GetByID(id string) (Endpoint, error)
	Create(endpoint Endpoint) (Endpoint, error)
	Update(endpoint Endpoint) (Endpoint, error)
	Delete(id string) error
	CreateDelivery(delivery Delivery) (Delivery, error)
	UpdateDelivery(delivery Delivery) (Delivery, error)
	GetDeliveriesByEndpointID(endpointID string) ([]Delivery, error)
	GetDeliveryByID(id string) (Delivery, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	_ = db.AutoMigrate(&Endpoint{}, &Delivery{})
	return &repo{db: db}
}

func (r *repo) GetAll() ([]Endpoint, error) {
	var endpoints []Endpoint
	if err := r.db.Order("created_at desc").Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (r *repo) GetActive() ([]Endpoint, error) {
	var endpoints []Endpoint
	if err := r.db.Where("is_active = ?", true).Order("created_at desc").Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (r *repo) GetByUserID(userID string) ([]Endpoint, error) {
	var endpoints []Endpoint
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&endpoints).Error; err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (r *repo) GetByID(id string) (Endpoint, error) {
	var endpoint Endpoint
	if err := r.db.First(&endpoint, "id = ?", id).Error; err != nil {
		return Endpoint{}, err
	}
	return endpoint, nil
}

func (r *repo) Create(endpoint Endpoint) (Endpoint, error) {
	if err := r.db.Create(&endpoint).Error; err != nil {
		return Endpoint{}, err
	}
	return endpoint, nil
}

func (r *repo) Update(endpoint Endpoint) (Endpoint, error) {
	if err := r.db.Save(&endpoint).Error; err != nil {
		return Endpoint{}, err
	}
	return endpoint, nil
}

func (r *repo) Delete(id string) error {
	if err := r.db.Delete(&Endpoint{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *repo) CreateDelivery(delivery Delivery) (Delivery, error) {
	if err := r.db.Create(&delivery).Error; err != nil {
		return Delivery{}, err
	}
	return delivery, nil
}

func (r *repo) UpdateDelivery(delivery Delivery) (Delivery, error) {
	if err := r.db.Save(&delivery).Error; err != nil {
		return Delivery{}, err
	}
	return delivery, nil
}

func (r *repo) GetDeliveriesByEndpointID(endpointID string) ([]Delivery, error) {
	var deliveries []Delivery
	if err := r.db.Where("endpoint_id = ?", endpointID).Order("created_at desc").Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *repo) GetDeliveryByID(id string) (Delivery, error) {
	var delivery Delivery
	if err := r.db.First(&delivery, "id = ?", id).Error; err != nil {
		return Delivery{}, err
	}
	return delivery, nil
}
