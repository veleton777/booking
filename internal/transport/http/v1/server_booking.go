package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/veleton777/booking_api/internal/dto"
	"github.com/veleton777/booking_api/internal/pkg/httputil"
)

const oneDay = time.Hour * 24

type BookingServer struct {
	bookingSvc BookingSvc
	validator  *validator.Validate
	l          *zerolog.Logger
}

//go:generate mockery --name BookingSvc
type BookingSvc interface {
	Booking(ctx context.Context, req dto.BookingReq) error
}

func NewBookingSvcServer(bookingSvc BookingSvc, l *zerolog.Logger) *BookingServer {
	return &BookingServer{
		bookingSvc: bookingSvc,
		validator:  validator.New(),
		l:          l,
	}
}

// Booking godoc
//
//	@Summary		Create new booking
//	@Description	Create new booking
//	@Tags			booking
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.BookingReq	true	"BookingReqDTO"
//	@Success		201
//	@Failure		400		{object}  httputil.HTTPError
//	@Router			/v1/booking [post]
func (s *BookingServer) Booking(w http.ResponseWriter, r *http.Request) {
	var bookingReq dto.BookingReq
	if err := json.NewDecoder(r.Body).Decode(&bookingReq); err != nil {
		s.l.Warn().Err(err).Msg("booking parse json")
		httputil.NewBadRequestErr(w, "invalid json body format") //nolint:wrapcheck

		return
	}

	if err := s.validator.Struct(bookingReq); err != nil {
		s.l.Warn().Err(err).Msg("booking validate req")
		httputil.NewBadRequestErr(w, err.Error()) //nolint:wrapcheck

		return
	}

	if time.Time(bookingReq.From).After(time.Time(bookingReq.To)) ||
		time.Time(bookingReq.From).Before(time.Now().Add(oneDay)) {
		s.l.Warn().Msg("invalid dates range")
		httputil.NewBadRequestErr(w, "invalid dates range") //nolint:wrapcheck

		return
	}

	if err := s.bookingSvc.Booking(r.Context(), bookingReq); err != nil {
		if bCode, ok := errToBusinessCode(err); ok {
			httputil.NewBusinessErr(w, int(bCode)) //nolint:wrapcheck

			return
		}

		s.l.Err(err).Msg("booking api error")
		httputil.NewInternalServerErr(w) //nolint:wrapcheck

		return
	}

	httputil.NewCreatedResponse(w) //nolint:wrapcheck

	return
}
