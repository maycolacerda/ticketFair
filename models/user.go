package models

type User struct {
	UserID   string `json:"user_id" gorm:"primaryKey;uniqueIndex"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"password" gorm:"not null"`
	Username string `json:"username" gorm:"uniqueIndex"`
	// Unique identifier for the user
}
