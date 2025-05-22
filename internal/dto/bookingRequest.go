package dto

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type CreateBookingRequest struct {
	UserID  int64         `json:"user_id" validate:"reuired"`
	EventID int64         `json:"event_id" validate:"reuired"`
	SeatID  int64         `json:"seat_id" validate:"reuired"`
	Status  BookingStatus `json:"status" validate:"reuired,oneof=pending confirmed cancelled"`
}
