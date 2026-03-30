// models/address.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	AddressID string         `json:"address_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProfileID string         `json:"profile_id" gorm:"uniqueIndex;not null"`
	Street    string         `json:"street"     gorm:"not null"`
	City      string         `json:"city"       gorm:"not null"`
	State     string         `json:"state"      gorm:"not null"`
	Country   string         `json:"country"    gorm:"not null;default:'BR'"`
	ZipCode   string         `json:"zip_code"   gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
