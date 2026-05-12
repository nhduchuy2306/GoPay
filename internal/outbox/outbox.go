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
	CreatedAt   time.Time      `json:"created_at"`
}
