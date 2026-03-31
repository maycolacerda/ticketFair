// models/password_reset_models.go
package models

import (
	"time"
)

type PasswordReset struct {
	ResetID   string     `json:"reset_id"  gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string     `json:"user_id"   gorm:"not null;index"`
	Code      string     `json:"code"      gorm:"not null"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at"   gorm:"default:null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationship
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
