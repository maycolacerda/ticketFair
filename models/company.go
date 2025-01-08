package models

import "gorm.io/gorm"

// The `Company` struct defines the structure of a company entity with various fields such as ID, name,
// address, phone, email, and registry.
// @property  - 1. CompanyID: Unique identifier for the company
// @property {int} CompanyID - CompanyID is an integer field representing the unique identifier for a
// company.
// @property {string} CompanyName - The `CompanyName` property in the `Company` struct represents the
// name of the company. It is a string type field and is tagged with `json:"company_name"` for JSON
// marshaling and unmarshaling.
// @property {string} CompanyAddress - The `CompanyAddress` property in the `Company` struct represents
// the physical address of the company. It typically includes details such as street address, city,
// state, and postal code.
// @property {string} CompanyPhone - The `CompanyPhone` property in the `Company` struct represents the
// phone number of the company. It is of type `string` and is tagged with `json:"company_phone"` for
// JSON marshaling and unmarshaling.
// @property {string} CompanyEmail - CompanyEmail is a field in the Company struct that represents the
// email address of the company. It is tagged with `json:"company_email"` for JSON marshaling and
// unmarshaling.
// @property {string} CompanyRegistry - The `CompanyRegistry` property in the `Company` struct
// represents the registry number or identifier of the company. This could be a unique registration
// number assigned to the company by a government authority or regulatory body for identification and
// legal purposes. It is commonly used to uniquely identify a company in official records and documents
type Company struct {
	gorm.Model
	CompanyID       int    `json:"company_id" gorm:"primary_key"`
	CompanyName     string `json:"company_name"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	CompanyEmail    string `json:"company_email"`
	CompanyRegistry string `json:"company_registry"`
}
