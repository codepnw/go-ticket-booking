package domain

import "time"

type Booking struct {
	ID         int64
	UserID     int64
	EventID    int64
	SeatID     int64
	Status     string
	BookedAt   time.Time
	CanceledAt time.Time
}
