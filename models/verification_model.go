// models/verification.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Verification struct {
	VerificationID string         `json:"verification_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string         `json:"user_id"         gorm:"not null;index"`
	Type           string         `json:"type"            gorm:"not null"` // "email" | "phone"
	Code           string         `json:"code"            gorm:"not null"`
	ExpiresAt      time.Time      `json:"expires_at"      gorm:"not null"`
	UsedAt         *time.Time     `json:"used_at"         gorm:"default:null"`
	CreatedAt      time.Time      `json:"created_at"      gorm:"autoCreateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationship
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
