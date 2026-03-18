// models/event.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	EventID     string         `json:"event_id"             gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MerchantID  string         `json:"merchant_id"          gorm:"not null;index"`
	Name        string         `json:"name"                 gorm:"not null"`
	Description string         `json:"description"          gorm:"type:text"`
	Location    string         `json:"location"             gorm:"not null"`
	StartTime   time.Time      `json:"start_time"           gorm:"not null"`
	EndTime     time.Time      `json:"end_time"             gorm:"not null"`
	Capacity    int            `json:"capacity"             gorm:"not null;default:0"`
	Active      bool           `json:"active"               gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at"           gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Merchant Merchant `json:"merchant,omitempty" gorm:"foreignKey:MerchantID"`
}
