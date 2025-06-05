package dto

type CreateSeatRequest struct {
	SectionID  int64  `json:"section_id" validate:"required"`
	RowLabel   string `json:"row_label" validate:"required"`
	SeatNumber int    `json:"seat_number" validate:"required"`
}

type CreateSeatsRequest struct {
	Seats []*CreateSeatRequest `json:"seats" validate:"required,dive,required"`
}

type UpdateSeatRequest struct {
	SectionID   *int64  `json:"section_id"`
	RowLabel    *string `json:"row_label"`
	SeatNumber  *int    `json:"seat_number"`
	IsAvailable *bool   `json:"is_available"`
}
