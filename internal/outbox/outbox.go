package outbox

import (
	"gopay/internal/common"
	"time"

	"gorm.io/datatypes"
)

type Outbox struct {
	common.BaseModel
	AggregateID string         `gorm:"type:uuid;not null;index" json:"aggregate_id"`
	EventType   string         `gorm:"not null;index" json:"event_type"`
	Payload     datatypes.JSON `gorm:"type:jsonb;not null" json:"payload"`
	Published   bool           `gorm:"default:false;index" json:"published"`
	Processing  bool           `gorm:"default:false;index" json:"processing"`
	Attempts    int            `gorm:"not null;default:0" json:"attempts"`
	NextRetryAt *time.Time     `json:"next_retry_at,omitempty"`
	PublishedAt *time.Time     `json:"published_at,omitempty"`
	LastError   string         `json:"last_error,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}
