package app

import (
	"context"
	"errors"
	"log"

	"gopay/internal/outbox"
	"gopay/internal/webhook"

	"gorm.io/gorm"
)

func StartBackgroundWorkers(ctx context.Context, db *gorm.DB) {
	outboxRepo := outbox.NewRepository(db)
	webhookRepo := webhook.NewRepository(db)
	webhookWorker := webhook.NewWorker()
	dispatcher := outbox.NewDispatcher(outboxRepo, webhookRepo, webhookWorker)

	go func() {
		if err := dispatcher.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Println("outbox dispatcher stopped with error:", err)
		}
	}()
}
