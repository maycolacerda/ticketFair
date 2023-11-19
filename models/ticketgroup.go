package models

type TicketGroup struct {
	TicketgroupID int     `json:"ticketgroup_id"`
	EventID       int     `json:"event_id"`
	GroupName     string  `json:"group_name"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	Description   string  `json:"description"`
}
