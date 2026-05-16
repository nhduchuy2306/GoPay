package outbox

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(event Outbox) (Outbox, error)
	ClaimPending(limit int) ([]Outbox, error)
	MarkPublished(id string) error
	MarkFailed(id string, cause error, retryAfter time.Duration) error
	CountPending() (int64, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	_ = db.AutoMigrate(&Outbox{})
	return &repo{db: db}
}

func (r *repo) Create(event Outbox) (Outbox, error) {
	if err := r.db.Create(&event).Error; err != nil {
		return Outbox{}, err
	}
	return event, nil
}

func (r *repo) ClaimPending(limit int) ([]Outbox, error) {
	if limit <= 0 {
		limit = 1
	}

	var claimed []Outbox
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var events []Outbox
		now := time.Now()
		staleCutoff := now.Add(-5 * time.Minute)
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("published = ? AND (processing = ? OR updated_at <= ?) AND (next_retry_at IS NULL OR next_retry_at <= ?)", false, false, staleCutoff, now).
			Order("created_at asc").
			Limit(limit).
			Find(&events).Error; err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		ids := make([]string, 0, len(events))
		for _, event := range events {
			ids = append(ids, event.ID)
		}

		if err := tx.Model(&Outbox{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"processing": true,
				"attempts":   gorm.Expr("attempts + ?", 1),
				"updated_at": now,
			}).Error; err != nil {
			return err
		}

		for i := range events {
			events[i].Processing = true
			events[i].Attempts++
		}
		claimed = events
		return nil
	})
	if err != nil {
		return nil, err
	}
	return claimed, nil
}

func (r *repo) MarkPublished(id string) error {
	now := time.Now()
	return r.db.Model(&Outbox{}).Where("id = ?", id).Updates(map[string]any{
		"published":     true,
		"processing":    false,
		"published_at":  &now,
		"next_retry_at": nil,
		"last_error":    "",
		"updated_at":    now,
	}).Error
}

func (r *repo) MarkFailed(id string, cause error, retryAfter time.Duration) error {
	if cause == nil {
		cause = errors.New("unknown outbox error")
	}
	now := time.Now()
	nextRetryAt := now.Add(retryAfter)
	return r.db.Model(&Outbox{}).Where("id = ?", id).Updates(map[string]any{
		"processing":    false,
		"next_retry_at": &nextRetryAt,
		"last_error":    cause.Error(),
		"updated_at":    now,
	}).Error
}

func (r *repo) CountPending() (int64, error) {
	var count int64
	if err := r.db.Model(&Outbox{}).Where("published = ?", false).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
