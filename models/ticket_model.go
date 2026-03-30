// models/ticket.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	TicketID      string         `json:"ticket_id"      gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TransactionID string         `json:"transaction_id" gorm:"not null;index"`
	UserID        string         `json:"user_id"        gorm:"not null;index"`
	EventID       string         `json:"event_id"       gorm:"not null;index"`
	Status        string         `json:"status"         gorm:"not null;default:'active'"`
	CreatedAt     time.Time      `json:"created_at"     gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at"     gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	User  *User  `json:"user,omitempty"        gorm:"foreignKey:UserID"`
	Event *Event `json:"event,omitempty"       gorm:"foreignKey:EventID"`
}
