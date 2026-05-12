package pkg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDbConnection() *gorm.DB {
	dsn := "host=localhost user=postgres password=12345 dbname=go-pay port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
