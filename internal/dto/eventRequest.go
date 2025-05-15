package dto

import "time"

type EventRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	LocationID  int        `json:"location_id"`
}
