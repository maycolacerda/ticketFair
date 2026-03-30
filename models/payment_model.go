// models/payment.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment tracks every Stripe PaymentIntent lifecycle.
// A Payment becomes a Transaction+Ticket only after status = "succeeded".
type Payment struct {
	PaymentID       string `json:"payment_id"        gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          string `json:"user_id"           gorm:"not null;index"`
	EventID         string `json:"event_id"          gorm:"not null;index"`
	TransactionID   string `json:"transaction_id"    gorm:"index"` // set after success
	StripePaymentID string `json:"stripe_payment_id" gorm:"uniqueIndex;not null"`
	Amount          int64  `json:"amount"            gorm:"not null"` // cents
	Currency        string `json:"currency"          gorm:"not null;default:'brl'"`
	Quantity        int    `json:"quantity"          gorm:"not null;default:1"`
	Status          string `json:"status"            gorm:"not null;default:'pending'"`
	// pending | succeeded | failed | canceled | refunded
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	User  *User  `json:"user,omitempty"  gorm:"foreignKey:UserID"`
	Event *Event `json:"event,omitempty" gorm:"foreignKey:EventID"`
}
