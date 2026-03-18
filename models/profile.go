// models/profile.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ProfileID     string         `json:"profile_id"           gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        string         `json:"user_id"              gorm:"uniqueIndex;not null"`
	FirstName     string         `json:"first_name"           gorm:"not null"`
	LastName      string         `json:"last_name"            gorm:"not null"`
	PhoneNumber   string         `json:"phone_number"         gorm:"uniqueIndex;not null"`
	VerifiedEmail bool           `json:"verified_email"  gorm:"default:false"`
	VerifiedPhone bool           `json:"verified_phone"  gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at"           gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	User    User    `json:"user,omitempty"    gorm:"foreignKey:UserID"`
	Address Address `json:"address,omitempty" gorm:"foreignKey:ProfileID"` // ← replaces Address string
}
