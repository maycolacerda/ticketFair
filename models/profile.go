package models

type Profile struct {
	ProfileID   string `json:"profile_id" gorm:"primaryKey"`
	UserID      string `json:"user_id" gorm:"unique"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number" gorm:"unique"`
	Address     string `json:"address"`
}
