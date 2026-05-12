package transaction

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		log.Fatal(err)
	}
	return &repo{db: db}
}
