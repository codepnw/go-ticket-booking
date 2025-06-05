package dto

import "time"

type EventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	LocationID  int       `json:"location_id" validate:"gte=0"`
}

type EventUpdateRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	LocationID  *int       `json:"location_id"`
}

type LocationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Address     string `json:"address" validate:"required"`
	Capacity    int64  `json:"capacity" validate:"gte=0,lte=200"`
	OwnerID     int    `json:"owner_id" validate:"gte=0"`
}

type LocationUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Capacity    *int64  `json:"capacity"`
	OwnerID     *int    `json:"owner_id"`
}
