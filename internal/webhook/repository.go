package webhook

import "gorm.io/gorm"

type Repository interface {
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}
