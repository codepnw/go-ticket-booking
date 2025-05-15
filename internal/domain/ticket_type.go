package domain

import "time"

type TicketType struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
