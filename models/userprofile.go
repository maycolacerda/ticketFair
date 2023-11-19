package models

import "time"

type UserProfile struct {
	UserProfileID int
	UserID        int
	FirstName     string
	LastName      string
	Email         string
	BirthDay      time.Time
}
