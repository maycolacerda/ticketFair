package models

import "gorm.io/gorm"

// The `Ticket` struct represents a ticket with fields for ID, user ID, event ID, group, and code.
// @property  - 1. `TicketID`: An integer representing the unique identifier for the ticket.
// @property {int} TicketID - The `TicketID` field in the `Ticket` struct represents the unique
// identifier for a ticket. It is of type `int` and is tagged with `json:"ticket_id"` for JSON
// serialization.
// @property {int} UserID - The `UserID` property in the `Ticket` struct represents the ID of the user
// who owns the ticket.
// @property {int} EventID - The `EventID` property in the `Ticket` struct represents the unique
// identifier of the event for which the ticket is issued. It is an integer value and is tagged with
// `json:"event_id"` for JSON marshaling purposes.
// @property {int} TicketGroup - The `TicketGroup` property in the `Ticket` struct represents the group
// or category to which the ticket belongs. It is used to classify tickets into different groups based
// on certain criteria such as seating area, ticket type, or any other categorization needed for the
// tickets.
// @property {string} TicketCode - The `TicketCode` property in the `Ticket` struct represents a unique
// code assigned to a ticket. This code can be used for identification and validation purposes, such as
// for scanning at events or verifying ticket ownership.
type Ticket struct {
	gorm.Model
	TicketID    int    `json:"ticket_id"`
	UserID      int    `json:"user_id"`
	EventID     int    `json:"event_id"`
	TicketGroup int    `json:"ticket_group"`
	TicketCode  string `json:"ticket_code"`
}
