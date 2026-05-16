package webhook

import (
	"gopay/internal/common"
	"gopay/internal/user"
	"time"

	"gorm.io/datatypes"
)

type Endpoint struct {
	common.BaseModel
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	User      user.User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	URL       string    `gorm:"not null" json:"url"`
	Secret    string    `gorm:"not null" json:"-"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Delivery struct {
	common.BaseModel
	EndpointID   string         `gorm:"type:uuid;not null;index" json:"endpoint_id"`
	Endpoint     Endpoint       `gorm:"foreignKey:EndpointID;references:ID" json:"endpoint,omitempty"`
	EventID      string         `gorm:"type:uuid;not null;index" json:"event_id"`
	Status       string         `gorm:"not null" json:"status"`
	AttemptCount int            `gorm:"not null;default:0" json:"attempt_count"`
	NextRetryAt  *time.Time     `json:"next_retry_at,omitempty"`
	ResponseCode *int           `json:"response_code,omitempty"`
	ResponseBody datatypes.JSON `gorm:"type:jsonb" json:"response_body,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CreateEndpointRequest struct {
	URL      string `json:"url"`
	Secret   string `json:"secret,omitempty"`
	IsActive bool   `json:"is_active,omitempty"`
}

type UpdateEndpointRequest struct {
	URL      string `json:"url,omitempty"`
	Secret   string `json:"secret,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}
