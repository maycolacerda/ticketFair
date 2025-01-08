package models

import "gorm.io/gorm"

// The Event type represents an event with various properties such as event ID, title, description,
// start and end times, location, and company ID.
// @property  - 1. `EventID`: An integer representing the unique identifier for the event.
// @property {int} EventID - EventID is an integer that represents the unique identifier for an event.
// @property {string} EventTitle - The `EventTitle` property in the `Event` struct represents the title
// or name of the event. It is a string type field and is tagged with `json:"event_title"` for JSON
// marshaling and unmarshaling.
// @property {string} Description - The `Event` struct represents an event entity with the following
// properties:
// @property {string} StartTime - The `StartTime` property in the `Event` struct represents the time at
// which the event is scheduled to start. It is of type string and is tagged with `json:"start_time"`
// for JSON serialization.
// @property {string} EndTime - The `EndTime` property in the `Event` struct represents the time when
// the event is scheduled to end. It is of type string and is tagged with `json:"end_time"` for JSON
// serialization. This property typically stores the end time of the event in a specific format, such
// as "HH
// @property {string} Location - The `Location` property in the `Event` struct represents the physical
// location where the event is taking place. It could be a specific address, venue name, or any other
// relevant information about where the event will be held.
// @property {int} CompanyID - The `CompanyID` property in the `Event` struct represents the ID of the
// company associated with the event. It is an integer field that is tagged with `json:"company_id"`
// for JSON marshaling and unmarshaling purposes.
type Event struct {
	gorm.Model
	EventID     int    `json:"event_id"`
	EventTitle  string `json:"event_title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    string `json:"location"`
	CompanyID   int    `json:"company_id"`
}
