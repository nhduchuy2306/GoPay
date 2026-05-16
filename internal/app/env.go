package app

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Println("warning: cannot load .env:", err)
	}
}
