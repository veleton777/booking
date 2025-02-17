package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/veleton777/booking_api/internal/booking/v1"
	"github.com/veleton777/booking_api/internal/booking/v1/booking/entity"
	"github.com/veleton777/booking_api/internal/booking/v1/booking/storage/memory"
	"github.com/veleton777/booking_api/internal/config"
	"github.com/veleton777/booking_api/internal/shutdown"
	"github.com/veleton777/booking_api/internal/transport/http/v1"
	memoryevent "github.com/veleton777/booking_api/pkg/event/memory"
	memorytx "github.com/veleton777/booking_api/pkg/transaction/memory"
)

type API struct {
	config *config.Config
	sh     *shutdown.Shutdown
	l      *zerolog.Logger

	bookingServer *v1.BookingServer
	bookingSvc    *booking.Svc
}

func New(_ context.Context, config *config.Config, l *zerolog.Logger) (*API, error) {
	sh := shutdown.New()

	a := &API{ //nolint:exhaustruct
		config: config,
		sh:     sh,
		l:      l,
	}

	bookingStorage := memory.NewStorage()

	// init rooms availability
	date := func(y, m, d int) int64 {
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local).Unix()
	}

	// hotelID->RoomID->date
	//nolint:mnd,gofumpt
	bookingStorage.Init(map[string]map[string]map[int64]*entity.Booking{
		"hotel_1": {
			"room_1": {
				date(2025, 4, 12): nil,
				date(2025, 4, 13): nil,
				date(2025, 4, 14): nil,
			},
			"room_2": {
				date(2025, 5, 9):  nil,
				date(2025, 5, 10): nil,
			},
		},
	})

	memoryTx := memorytx.NewTxClient()
	eventSvc := memoryevent.NewEventClient()

	bookingSvc := booking.NewBookingSvc(bookingStorage, memoryTx, eventSvc, l)
	a.bookingSvc = bookingSvc

	a.bookingServer = v1.NewBookingSvcServer(bookingSvc, l)

	return a, nil
}

func (s *API) Run(ctx context.Context) error {
	var err error

	r := chi.NewRouter()

	// Middleware
	r.Use(LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	s.routes(r)

	server := &http.Server{
		Addr:              s.config.HTTPAddr(),
		Handler:           r,
		ReadHeaderTimeout: time.Minute,
	}

	s.sh.AddNormalPriority(func(_ context.Context) error {
		if err := server.Close(); err != nil {
			return errors.Wrap(err, "http server shutdown")
		}

		return nil
	})

	go func() {
		s.l.Info().Msg("start server")

		if serverErr := server.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			err = errors.Wrap(serverErr, "http server listen port")
		}
	}()

	s.sh.WaitShutdown(ctx)

	return err
}
