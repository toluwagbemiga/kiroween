package repository

import "errors"

var (
	// ErrNotFound is returned when a record is not found
	ErrNotFound = errors.New("record not found")
)
