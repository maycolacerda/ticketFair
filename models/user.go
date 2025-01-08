package models

import "gorm.io/gorm"

// The `User` struct defines a user model with fields for user ID, username, and password.
// @property  - The `User` struct represents a user entity in a Go program. Here are the properties of
// the `User` struct:
// @property {int} UserID - The `UserID` field in the `User` struct represents the unique identifier
// for each user in the system. It is of type `int` and is tagged with `json:"user_id"` for JSON
// serialization.
// @property {string} Username - The `Username` property in the `User` struct represents the username
// of a user. It is a string type field and is tagged with `json:"username"` for JSON serialization
// purposes.
// @property {string} Password - The `Password` property in the `User` struct represents the password
// of a user account. It is typically used for authentication purposes to verify the identity of a user
// when they log in to a system or application. It is important to securely store and handle passwords
// to protect user data and ensure the security
type User struct {
	gorm.Model
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
