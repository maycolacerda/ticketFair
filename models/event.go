package models

type Event struct {
	EventID     string `gorm:"primaryKey" json:"event_id"`
	MerchantID  string `json:"merchant_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Date        string `json:"date" binding:"required"`
	Location    string `json:"location" binding:"required"`
}
