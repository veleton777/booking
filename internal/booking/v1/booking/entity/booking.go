package entity

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID
	HotelID   string
	RoomID    string
	UserEmail string
	From      time.Time
	To        time.Time
	CreatedAt time.Time
}
