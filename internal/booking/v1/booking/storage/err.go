package storage

import "errors"

var (
	ErrPlaceNotAvailable = errors.New("place not available")
	ErrEntityNotFound    = errors.New("entity not found")
)
