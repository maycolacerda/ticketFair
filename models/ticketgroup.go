package models

type TicketGroup struct {
	ID        int
	EventID   int
	GroupName string
	Price     float64
	Quantity  int
}
