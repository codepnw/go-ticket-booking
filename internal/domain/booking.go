package domain

import "time"

type Booking struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	EventID     int64      `json:"event_id"`
	SeatID      int64      `json:"seat_id"`
	Status      string     `json:"status"`
	BookedAt    time.Time  `json:"booked_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
}
