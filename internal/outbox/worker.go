package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"gopay/internal/utils/enum"
	"gopay/internal/webhook"

	"gorm.io/datatypes"
)

type Dispatcher struct {
	repo           Repository
	webhookRepo    webhook.Repository
	webhookWorker  webhook.Worker
	batchSize      int
	concurrency    int
	endpointFanout int
	pollInterval   time.Duration
	logger         *log.Logger
}

type TransactionCompletedPayload struct {
	TransactionID string `json:"transaction_id"`
	Amount        int64  `json:"amount"`
	FromWalletID  string `json:"from_wallet_id"`
	ToWalletID    string `json:"to_wallet_id"`
	Status        string `json:"status"`
	Description   string `json:"description,omitempty"`
}

func NewDispatcher(repo Repository, webhookRepo webhook.Repository, webhookWorker webhook.Worker) *Dispatcher {
	return &Dispatcher{
		repo:           repo,
		webhookRepo:    webhookRepo,
		webhookWorker:  webhookWorker,
		batchSize:      20,
		concurrency:    4,
		endpointFanout: 4,
		pollInterval:   2 * time.Second,
		logger:         log.Default(),
	}
}

func (d *Dispatcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(d.pollInterval)
	defer ticker.Stop()

	for {
		if err := d.ProcessOnce(ctx); err != nil && !errors.Is(err, context.Canceled) {
			d.logger.Println("outbox process error:", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (d *Dispatcher) ProcessOnce(ctx context.Context) error {
	if ctx != nil && ctx.Err() != nil {
		return ctx.Err()
	}

	batch, err := d.repo.ClaimPending(d.batchSize)
	if err != nil {
		return err
	}
	if len(batch) == 0 {
		return nil
	}

	sem := make(chan struct{}, d.concurrency)
	var wg sync.WaitGroup
	for _, event := range batch {
		event := event
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			if err := d.processEvent(event); err != nil {
				d.logger.Printf("outbox event %s failed: %v", event.ID, err)
			}
		}()
	}
	wg.Wait()
	return nil
}

func (d *Dispatcher) processEvent(event Outbox) error {
	var payload TransactionCompletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return d.repo.MarkFailed(event.ID, err, retryDelay(event.Attempts))
	}

	if event.EventType != enum.TransactionCompleted {
		return d.repo.MarkFailed(event.ID, errors.New("unsupported outbox event type"), retryDelay(event.Attempts))
	}

	endpoints, err := d.webhookRepo.GetActive()
	if err != nil {
		return d.repo.MarkFailed(event.ID, err, retryDelay(event.Attempts))
	}

	if len(endpoints) == 0 {
		return d.repo.MarkPublished(event.ID)
	}

	endpointSem := make(chan struct{}, d.endpointFanout)
	var endpointWG sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return d.repo.MarkFailed(event.ID, err, retryDelay(event.Attempts))
	}

	for _, endpoint := range endpoints {
		endpoint := endpoint
		endpointSem <- struct{}{}
		endpointWG.Add(1)
		go func() {
			defer endpointWG.Done()
			defer func() { <-endpointSem }()
			if err := d.deliverToEndpoint(event, endpoint, payloadBytes); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}()
	}
	endpointWG.Wait()

	if firstErr != nil {
		return d.repo.MarkFailed(event.ID, firstErr, retryDelay(event.Attempts))
	}
	return d.repo.MarkPublished(event.ID)
}

func (d *Dispatcher) deliverToEndpoint(event Outbox, endpoint webhook.Endpoint, payloadBytes []byte) error {
	delivery := webhook.Delivery{
		EndpointID:   endpoint.ID,
		EventID:      event.ID,
		Status:       enum.StatusPending,
		AttemptCount: 1,
	}
	created, err := d.webhookRepo.CreateDelivery(delivery)
	if err != nil {
		return err
	}

	statusCode, responseBody, err := d.webhookWorker.Deliver(endpoint, payloadBytes)
	if err != nil {
		created.Status = enum.StatusFailed
		created.ResponseCode = &statusCode
		created.ResponseBody = datatypes.JSON(responseBody)
		_, updateErr := d.webhookRepo.UpdateDelivery(created)
		if updateErr != nil {
			return updateErr
		}
		return err
	}

	created.Status = enum.StatusSuccess
	created.ResponseCode = &statusCode
	created.ResponseBody = datatypes.JSON(responseBody)
	_, err = d.webhookRepo.UpdateDelivery(created)
	return err
}

func retryDelay(attempts int) time.Duration {
	if attempts <= 0 {
		return 5 * time.Second
	}
	backoff := time.Duration(1<<minInt(attempts, 6)) * time.Second
	if backoff > 5*time.Minute {
		return 5 * time.Minute
	}
	return backoff
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var _ = NewDispatcher
