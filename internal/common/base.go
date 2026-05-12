package common

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (b *BaseModel) BeforeCreate() {
	b.ID = uuid.NewString()
	b.CreatedAt = time.Now()
}
