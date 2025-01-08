package models

import "gorm.io/gorm"

// The `TicketGroup` struct represents a group of tickets for an event with various attributes such as
// group name, price, quantity, and description.
// @property  - 1. `TicketgroupID`: An integer representing the ID of the ticket group.
// @property {int} TicketgroupID - TicketgroupID is an integer that represents the unique identifier
// for a ticket group.
// @property {int} EventID - EventID is the identifier for the event associated with the ticket group.
// @property {string} GroupName - The `GroupName` property in the `TicketGroup` struct represents the
// name of the ticket group. It is a string type field that holds the name or title of the ticket
// group.
// @property {float64} Price - The `Price` property in the `TicketGroup` struct represents the price of
// a ticket group for an event. It is of type `float64` and is tagged with `json:"price"` for JSON
// serialization.
// @property {int} Quantity - The `Quantity` property in the `TicketGroup` struct represents the number
// of tickets available in that ticket group. It indicates how many tickets can be purchased for that
// specific group.
// @property {string} Description - The `TicketGroup` struct represents a group of tickets for an
// event. Here is a description of each property:
// @property {string} GroupCode - The `GroupCode` property in the `TicketGroup` struct represents a
// unique code assigned to a ticket group. This code can be used for identification or categorization
// purposes.
type TicketGroup struct {
	gorm.Model
	TicketgroupID    int     `json:"ticketgroup_id"`
	EventID          int     `json:"event_id"`
	GroupName        string  `json:"group_name"`
	GroupPrice       float64 `json:"price"`
	GroupQuantity    int     `json:"quantity"`
	GroupDescription string  `json:"description"`
	GroupCode        string  `json:"group_code"`
}
