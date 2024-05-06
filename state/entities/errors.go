package entities

import "errors"

var (
	// ErrNotFound indicates the state database may be empty
	ErrNotFound = errors.New("not found")
	// ErrStorageNotFound is used when the object is not found in the storage
	ErrStorageNotFound = errors.New("not found in the Storage")
	// ErrStorageNotRegister is used when the object is not found in the synchronizer
	ErrStorageNotRegister = errors.New("not registered storage")
)
