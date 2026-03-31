// models/ticket_type.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketType struct {
	TicketTypeID string         `json:"ticket_type_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EventID      string         `json:"event_id"       gorm:"not null;index"`
	Name         string         `json:"name"           gorm:"not null"`
	Description  string         `json:"description"    gorm:"type:text"`
	Category     string         `json:"category"       gorm:"not null;default:'general'"`
	PriceCents   int64          `json:"price_cents"    gorm:"not null"`
	Capacity     int            `json:"capacity"       gorm:"not null"`
	Available    int            `json:"available"      gorm:"not null"`
	MinPerOrder  int            `json:"min_per_order"  gorm:"not null;default:1"`
	MaxPerOrder  int            `json:"max_per_order"  gorm:"not null;default:10"`
	SaleStartsAt *time.Time     `json:"sale_starts_at" gorm:"default:null"`
	SaleEndsAt   *time.Time     `json:"sale_ends_at"   gorm:"default:null"`
	Active       bool           `json:"active"         gorm:"default:true"`
	SortOrder    int            `json:"sort_order"     gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"     gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at"     gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID"`
}
