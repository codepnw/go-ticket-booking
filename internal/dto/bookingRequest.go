package dto

import "time"

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type CreateBookingRequest struct {
	UserID  int64 `json:"user_id" validate:"required"`
	EventID int64 `json:"event_id" validate:"required"`
	SeatID  int64 `json:"seat_id" validate:"required"`
}

type UpdateBookingRequest struct {
	UserID  *int64 `json:"user_id"`
	EventID *int64 `json:"event_id"`
	SeatID  *int64 `json:"seat_id"`
}

type BookingResponse struct {
	ID          int64        `json:"id"`
	User        bookingUser  `json:"user"`
	Event       bookingEvent `json:"event"`
	Seat        bookingSeat  `json:"seat"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	ConfirmedAt *time.Time   `json:"confirmed_at"`
	CancelledAt *time.Time   `json:"cancelled_at"`
}

type BookingSeatUpdateRequest struct {
	SeatID int64 `json:"seat_id" validate:"required"`
}

type bookingEvent struct {
	EventID   int64  `json:"event_id"`
	EventName string `json:"event_name"`
}

type bookingSeat struct {
	SeatID     int64  `json:"seat_id"`
	SeatNumber int    `json:"seat_number"`
	RowLabel   string `json:"row_label"`
}

type bookingUser struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
