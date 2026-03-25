package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	AdminID   string         `json:"admin_id"             gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name"                 gorm:"not null"`
	Email     string         `json:"email"                gorm:"uniqueIndex;not null"`
	Password  string         `json:"-"                    gorm:"not null"`
	Active    bool           `json:"active"               gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at"           gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
