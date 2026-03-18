// models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID    string         `json:"user_id"              gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string         `json:"email"                gorm:"uniqueIndex;not null"`
	Password  string         `json:"-"                    gorm:"not null"`
	Username  string         `json:"username"             gorm:"uniqueIndex;not null"`
	Active    bool           `json:"active"               gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at"           gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
