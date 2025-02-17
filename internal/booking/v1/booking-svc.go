package booking

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/veleton777/booking_api/internal/booking/v1/booking/entity"
	"github.com/veleton777/booking_api/internal/booking/v1/booking/storage"
	"github.com/veleton777/booking_api/internal/dto"
	"github.com/veleton777/booking_api/pkg/event"
	"github.com/veleton777/booking_api/pkg/transaction"
)

type Svc struct {
	bookingStorage Storage
	txSvc          transaction.Tx
	eventSvc       event.Event
	l              *zerolog.Logger
}

func NewBookingSvc(
	bookingStorage Storage,
	txSvc transaction.Tx,
	eventSvc event.Event,
	l *zerolog.Logger,
) *Svc {
	return &Svc{
		bookingStorage: bookingStorage,
		txSvc:          txSvc,
		eventSvc:       eventSvc,
		l:              l,
	}
}

func (s *Svc) Booking(ctx context.Context, req dto.BookingReq) error {
	ctx, err := s.txSvc.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	defer func() {
		if err != nil {
			if rollbackErr := s.txSvc.Rollback(ctx); rollbackErr != nil {
				s.l.Err(rollbackErr).Msg("try rollback tx")
			}
		}
	}()

	bookingEnt := entity.Booking{
		ID:        uuid.New(),
		HotelID:   req.HotelID,
		RoomID:    req.RoomID,
		UserEmail: req.UserEmail,
		From:      time.Time(req.From),
		To:        time.Time(req.To),
		CreatedAt: time.Now(),
	}

	if err = s.bookingStorage.SaveBooking(ctx, bookingEnt); err != nil {
		if errors.Is(err, storage.ErrPlaceNotAvailable) {
			return entity.ErrRoomNotAvailable
		}

		if errors.Is(err, storage.ErrEntityNotFound) {
			return entity.ErrEntityNotFound
		}

		return errors.Wrap(err, "save booking in storage")
	}

	b, err := json.Marshal(bookingEnt)
	if err != nil {
		return errors.Wrap(err, "marshal booking entity to json")
	}

	ev, err := event.NewEvent(uuid.New(), event.TypeCreated, string(b))
	if err != nil {
		return errors.Wrap(err, "create event entity")
	}

	if err = s.eventSvc.SaveEvent(ctx, ev); err != nil {
		return errors.Wrap(err, "save event")
	}

	if err = s.txSvc.Commit(ctx); err != nil {
		return errors.Wrap(err, "tx commit")
	}

	return nil
}
