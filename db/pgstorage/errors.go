package pgstorage

import "errors"

var (
	// ErrNotFound indicates the state database may be empty
	ErrNotFound = errors.New("not found in db")
)
