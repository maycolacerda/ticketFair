package models

import "gorm.io/gorm"

// The `CompanyUser` struct defines the structure of a user belonging to a company in Go.
// @property  - 1. `CompanyUserID`: Unique identifier for the company user.
// @property {int} CompanyUserID - CompanyUserID is an integer that represents the unique identifier
// for a company user.
// @property {int} CompanyID - The `CompanyID` property in the `CompanyUser` struct represents the
// unique identifier for the company to which the user belongs.
// @property {string} CompanyUserLogin - The `CompanyUserLogin` property in the `CompanyUser` struct
// represents the login username of a company user. It is used to uniquely identify the user when
// logging into the system.
// @property {string} CompanyUserPassword - The `CompanyUserPassword` property in the `CompanyUser`
// struct represents the password of a company user. It is a string type field used to store the
// password associated with a company user account.
// @property {string} ComanyUserRole - It looks like there is a typo in the struct definition. The
// property `ComanyUserRole` should be `CompanyUserRole`.
type CompanyUser struct {
	gorm.Model
	CompanyUserID       int    `json:"company_user_id"`
	CompanyID           int    `json:"company_id"`
	CompanyUserLogin    string `json:"company_user_login"`
	CompanyUserPassword string `json:"company_user_password"`
	ComanyUserRole      string `json:"company_user_role"`
}
