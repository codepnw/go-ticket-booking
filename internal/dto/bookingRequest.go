package dto

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusBooked    BookingStatus = "booked"
	StatusCancelled BookingStatus = "cancelled"
)

type CreateBookingRequest struct {
	UserID  int64 `json:"user_id" validate:"required"`
	EventID int64 `json:"event_id" validate:"required"`
	SeatID  int64 `json:"seat_id" validate:"required"`
	// Status  BookingStatus `json:"status" validate:"reuired,oneof=pending confirmed cancelled"`
}

type UpdateBookingRequest struct {
	UserID  *int64 `json:"user_id"`
	EventID *int64 `json:"event_id"`
	SeatID  *int64 `json:"seat_id"`
}
