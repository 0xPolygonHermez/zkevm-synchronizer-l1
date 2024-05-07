package entities

import "errors"

var (
	// ErrNotFound indicates the state database may be empty
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists indicates that this entity already exists
	ErrAlreadyExists = errors.New("already exists")
	// ErrForeignKeyViolation indicates that a foreign key constraint was violated
	ErrForeignKeyViolation = errors.New("foreign key violation")
	// ErrStorageNotFound is used when the object is not found in the storage
	ErrStorageNotFound = errors.New("not found in the Storage")
	// ErrStorageNotRegister is used when the object is not found in the synchronizer
	ErrStorageNotRegister = errors.New("not registered storage")
)
