package dto

import "time"

type SectionRequest struct {
	EventID   int64  `json:"event_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	SeatCount int    `json:"seat_count" validate:"gte=0,lte=200"`
}

type SectionUpdate struct {
	EventID   *int64    `json:"event_id"`
	Name      *string   `json:"name"`
	SeatCount *int      `json:"seat_count"`
	UpdatedAt time.Time `json:"updated_at"`
}
