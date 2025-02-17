package memory

import (
	"context"
	"sync"
	"time"

	"github.com/veleton777/booking_api/internal/booking/v1/booking/entity"
	"github.com/veleton777/booking_api/internal/booking/v1/booking/storage"
)

type Storage struct {
	// hotelID->RoomID->date
	storage map[string]map[string]map[int64]*entity.Booking

	// todo btree index for search user bookings
	mu *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]map[string]map[int64]*entity.Booking),
		mu:      &sync.Mutex{},
	}
}

func (s *Storage) Init(m map[string]map[string]map[int64]*entity.Booking) {
	s.storage = m
}

//nolint:cyclop
func (s *Storage) SaveBooking(_ context.Context, booking entity.Booking) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.storage[booking.HotelID][booking.RoomID]; !ok {
		return storage.ErrEntityNotFound
	}

	if len(s.storage[booking.HotelID]) == 0 {
		s.storage[booking.HotelID] = make(map[string]map[int64]*entity.Booking)
	}

	if len(s.storage[booking.HotelID][booking.RoomID]) == 0 {
		s.storage[booking.HotelID][booking.RoomID] = make(map[int64]*entity.Booking)
	}

	dates := getDatesBetween(booking.From, booking.To)

	var (
		isAlreadyBookedForUser bool
		isAllowDate            bool
	)

	for _, d := range dates {
		b, ok := s.storage[booking.HotelID][booking.RoomID][d.Unix()]
		if !ok {
			return storage.ErrPlaceNotAvailable
		}

		if b == nil {
			isAllowDate = true

			continue
		}

		if b.UserEmail != booking.UserEmail {
			return storage.ErrPlaceNotAvailable
		}

		isAlreadyBookedForUser = true
	}

	if isAlreadyBookedForUser {
		if isAllowDate {
			return storage.ErrPlaceNotAvailable
		}

		// idempotency: user already booked these dates
		return nil
	}

	for _, d := range dates {
		s.storage[booking.HotelID][booking.RoomID][d.Unix()] = &booking
	}

	return nil
}

func getDatesBetween(from, to time.Time) []time.Time {
	var dates []time.Time

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}

	return dates
}
