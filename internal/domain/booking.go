package domain

import "time"

type Booking struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	EventID     int64      `json:"event_id"`
	SeatID      int64      `json:"seat_id"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
}
