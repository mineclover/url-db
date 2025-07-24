package repository

import "errors"

// Common repository errors
var (
	// ErrNotFound is returned when a requested entity is not found
	ErrNotFound = errors.New("entity not found")

	// ErrDuplicateKey is returned when trying to create an entity with a duplicate key
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrConstraintViolation is returned when a database constraint is violated
	ErrConstraintViolation = errors.New("constraint violation")

	// ErrConcurrencyConflict is returned when a concurrency conflict occurs
	ErrConcurrencyConflict = errors.New("concurrency conflict")
)
