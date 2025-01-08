package models

import (
	"time"

	"gorm.io/gorm"
)

// The `UserProfile` struct represents a user profile with fields for ID, user ID, first name, last
// name, email, and birth day.
// @property  - The `UserProfile` struct represents a user profile entity in a Go program. Here are the
// properties of the struct:
// @property {int} UserProfileID - The `UserProfileID` field in the `UserProfile` struct represents the
// unique identifier for a user profile. It is of type `int` and is tagged with
// `json:"user_profile_id"` for JSON serialization purposes.
// @property {int} UserID - The `UserID` field in the `UserProfile` struct represents the unique
// identifier of the user associated with the profile.
// @property {string} FirstName - The `FirstName` property in the `UserProfile` struct represents the
// first name of a user. It is a string type field and is tagged with `json:"first_name"` for JSON
// marshaling purposes.
// @property {string} LastName - The `LastName` property in the `UserProfile` struct represents the
// last name of a user. It is a string type field that stores the last name of the user.
// @property {string} Email - The `UserProfile` struct represents a user profile in a system. The
// `Email` field in the struct is used to store the email address of the user.
// @property BirthDay - The `BirthDay` property in the `UserProfile` struct represents the birth date
// of the user. It is of type `time.Time` and is tagged with `json:"birth_day"` for JSON serialization.
type UserProfile struct {
	gorm.Model
	UserProfileID int       `json:"user_profile_id"`
	UserID        int       `json:"user_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	BirthDay      time.Time `json:"birth_day"`
}
