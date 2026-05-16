package main

import (
	"context"
	"gopay/internal/app"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.New(ctx)
	if err := a.Run(":8888"); err != nil {
		log.Fatal(err)
	}
}
