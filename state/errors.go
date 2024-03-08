package state

import "errors"

var (
	// ErrNotFound is used when the object is not found
	ErrNotFound = errors.New("not found")
)
