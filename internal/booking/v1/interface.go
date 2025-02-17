package booking

import (
	"context"

	"github.com/veleton777/booking_api/internal/booking/v1/booking/entity"
)

//go:generate mockery --name Storage
type Storage interface {
	SaveBooking(_ context.Context, booking entity.Booking) error
}
