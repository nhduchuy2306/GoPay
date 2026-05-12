package app

import (
	"log"

	"github.com/joho/godotenv"
)

func InitEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}
