// models/transaction.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	TransactionID string         `json:"transaction_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        string         `json:"user_id"        gorm:"not null;index"`
	EventID       string         `json:"event_id"       gorm:"not null;index"`
	Amount        float64        `json:"amount"         gorm:"not null"`
	TicketTypeID  string         `json:"ticket_type_id" gorm:"index"`
	Quantity      int            `json:"quantity"       gorm:"not null;default:1"`
	Status        string         `json:"status"         gorm:"not null;default:'pending'"`
	CreatedAt     time.Time      `json:"created_at"     gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at"     gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	User       *User       `json:"user,omitempty"  gorm:"foreignKey:UserID"`
	Event      *Event      `json:"event,omitempty" gorm:"foreignKey:EventID"`
	TicketType *TicketType `json:"ticket_type,omitempty" gorm:"foreignKey:TicketTypeID"`
}
