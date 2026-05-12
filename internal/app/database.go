package app

import (
	"gopay/pkg"
	"log"

	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := pkg.CreateDbConnection()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
