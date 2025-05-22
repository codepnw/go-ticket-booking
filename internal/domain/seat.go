package domain

type Seat struct {
	ID          int64  `json:"id"`
	SectionID   int64  `json:"section_id"`
	RowLabel    string `json:"row_label"`
	SeatNumber  int    `json:"seat_number"`
	IsAvailable bool   `json:"is_available"`
}
