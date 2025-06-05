package errs

import "errors"

var (
	ErrSeatNotFound     = errors.New("seat not found")
	ErrEventNotFound    = errors.New("event not found")
	ErrLocationNotFound = errors.New("location not found")

	ErrNoFieldsToUpdate = errors.New("no fields to update")
)
