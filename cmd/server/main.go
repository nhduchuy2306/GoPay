package main

import (
	"gopay/internal/app"
	"log"
)

func main() {
	a := app.New()
	if err := a.Run(":8888"); err != nil {
		log.Fatal(err)
	}
}
