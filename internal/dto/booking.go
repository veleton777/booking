package dto

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

//nolint:tagliatelle
type BookingReq struct {
	HotelID   string `json:"hotel_id" validate:"required" example:"123"`
	RoomID    string `json:"room_id" validate:"required" example:"456"`
	UserEmail string `json:"email" validate:"required,email" example:"user@gmail.com"`
	From      Date   `json:"from" validate:"required" example:"2025-04-20"`
	To        Date   `json:"to" validate:"required" example:"2025-04-25"`
}

type Date time.Time

func (t *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}

	var (
		parsed time.Time
		err    error
	)

	layouts := []string{
		"2006-01-02",
	}

	for _, layout := range layouts {
		parsed, err = time.Parse(layout, s)
		if err == nil {
			break
		}
	}

	if err != nil {
		return errors.Wrap(err, "wrong time format")
	}

	*t = Date(parsed.UTC())

	return nil
}

type BookingResp struct{}
