package entity

import "errors"

var (
	ErrEntityNotFound   = errors.New("entity not found")
	ErrRoomNotAvailable = errors.New("room not available")
)
