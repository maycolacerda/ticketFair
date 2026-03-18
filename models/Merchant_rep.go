// models/merchantRep.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type MerchantRep struct {
	MerchantRepID string         `json:"merchant_rep_id"      gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MerchantID    string         `json:"merchant_id"          gorm:"not null;index"`
	Name          string         `json:"name"                 gorm:"not null"`
	Email         string         `json:"email"                gorm:"uniqueIndex;not null"`
	Phone         string         `json:"phone"                gorm:"not null"`
	Password      string         `json:"-"                    gorm:"not null"`
	Role          string         `json:"role"                 gorm:"not null;default:'staff'"`
	Active        bool           `json:"active"               gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at"           gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationship
	Merchant Merchant `json:"merchant,omitempty" gorm:"foreignKey:MerchantID"`
}
