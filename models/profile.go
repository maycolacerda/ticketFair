package models

import "time"

type Profile struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	BirthDay  time.Time
}
