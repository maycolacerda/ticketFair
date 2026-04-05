// models/connection.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Connection struct {
	ConnectionID string         `json:"connection_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RequesterID  string         `json:"requester_id"  gorm:"not null;index"`
	AddresseeID  string         `json:"addressee_id"  gorm:"not null;index"`
	Status       string         `json:"status"        gorm:"not null;default:'pending'"`
	CreatedAt    time.Time      `json:"created_at"    gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at"    gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Requester *User `json:"requester,omitempty" gorm:"foreignKey:RequesterID"`
	Addressee *User `json:"addressee,omitempty" gorm:"foreignKey:AddresseeID"`
}
