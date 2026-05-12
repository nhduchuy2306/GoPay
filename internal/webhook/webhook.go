package webhook

import (
	"gopay/internal/common"
	"gopay/internal/user"
	"time"

	"gorm.io/datatypes"
)

type WebhookDeliveryStatus string

const (
	Pending WebhookDeliveryStatus = "PENDING"
	Success WebhookDeliveryStatus = "SUCCESS"
	Failed  WebhookDeliveryStatus = "FAILED"
)

type WebhookEndpoint struct {
	common.BaseModel
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	User      user.User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	URL       string    `gorm:"not null" json:"url"`
	Secret    string    `gorm:"not null" json:"-"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookDelivery struct {
	common.BaseModel
	EndpointID   string                `gorm:"type:uuid;not null;index" json:"endpoint_id"`
	Endpoint     WebhookEndpoint       `gorm:"foreignKey:EndpointID;references:ID" json:"endpoint,omitempty"`
	EventID      string                `gorm:"type:uuid;not null;index" json:"event_id"`
	Status       WebhookDeliveryStatus `gorm:"not null" json:"status"`
	AttemptCount int                   `gorm:"not null;default:0" json:"attempt_count"`
	NextRetryAt  *time.Time            `json:"next_retry_at,omitempty"`
	ResponseCode *int                  `json:"response_code,omitempty"`
	ResponseBody datatypes.JSON        `gorm:"type:jsonb" json:"response_body,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}
