package errs

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrSeatNotFound     = errors.New("seat not found")
	ErrEventNotFound    = errors.New("event not found")
	ErrLocationNotFound = errors.New("location not found")

	ErrNoFieldsToUpdate = errors.New("no fields to update")
	ErrInvalidInputData = errors.New("invalid input data")
)
