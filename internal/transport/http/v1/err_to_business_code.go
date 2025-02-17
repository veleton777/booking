package v1

import (
	"errors"

	"github.com/veleton777/booking_api/internal/booking/v1/booking/entity"
)

type BusinessCode int

const (
	DefaultBusinessCode BusinessCode = 0
	RoomNotAvailable    BusinessCode = 1
	EntityNotFound      BusinessCode = 2
)

func errToBusinessCode(err error) (BusinessCode, bool) {
	if errors.Is(err, entity.ErrRoomNotAvailable) {
		return RoomNotAvailable, true
	}

	if errors.Is(err, entity.ErrEntityNotFound) {
		return EntityNotFound, true
	}

	return DefaultBusinessCode, false
}
