package dto

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

type UpdateBookingStatus struct {
	ID     int64         `json:"id" validate:"required"`
	Status BookingStatus `json:"status" validate:"required,oneof=confirmed cancelled"`
}