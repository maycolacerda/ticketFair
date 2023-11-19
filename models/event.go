package models

type Event struct {
	EventID     int    `json:"event_id"`
	EventTitle  string `json:"event_title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    string `json:"location"`
}
