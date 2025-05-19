package domain

import "time"

type Section struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"event_id"`
	Name      string    `json:"name"`
	SeatCount int       `json:"seat_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}	

type Seat struct {
	ID          int64  `json:"id"`
	SectionID   int64  `json:"section_id"`
	RowLabel    string `json:"row_label"`
	SeatNumber  int    `json:"seat_number"`
	IsAvailable bool   `json:"is_available"`
}
