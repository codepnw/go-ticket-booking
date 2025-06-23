package errs

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrSeatNotFound     = errors.New("seat not found")
	ErrSectionNotFound  = errors.New("section not found")
	ErrEventNotFound    = errors.New("event not found")
	ErrLocationNotFound = errors.New("location not found")
	ErrBookingNotFound  = errors.New("booking not found")

	ErrNoFieldsToUpdate        = errors.New("no fields to update")
	ErrInvalidInputData        = errors.New("invalid input data")
	ErrInvalidSeatEvent        = errors.New("invalid seat event")
	ErrSeatAlreadyBooked       = errors.New("seat already booked")
	ErrBookingAlreadyConfirmed = errors.New("booking already confirmed")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
	ErrBookingNotPending       = errors.New("cannot update status confirmed or cancelled booking")
)
