package models

type Ticket struct {
	TicketID    int `json:"ticket_id"`
	UserID      int `json:"user_id"`
	EventID     int `json:"event_id"`
	TicketGroup int `json:"ticket_group"`
}
