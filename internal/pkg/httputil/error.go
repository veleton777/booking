package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type HTTPError struct {
	Code         int    `json:"code"`
	Text         string `json:"text"`
	BusinessCode int    `json:"businessCode,omitempty"`
}

func NewNotFoundErr(w http.ResponseWriter) error {
	resp := HTTPError{
		Code:         http.StatusNotFound,
		Text:         "Not Found",
		BusinessCode: 0,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return errors.Wrap(err, "write json msg resp")
	}

	return nil
}

func NewBadRequestErr(w http.ResponseWriter, msg string) error {
	resp := HTTPError{
		Code:         http.StatusBadRequest,
		Text:         "Bad Request",
		BusinessCode: 0,
	}

	if msg != "" {
		resp.Text = msg
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return errors.Wrap(err, "write json msg resp")
	}

	return nil
}

func NewBusinessErr(w http.ResponseWriter, businessCode int) error {
	resp := HTTPError{
		Code:         http.StatusBadRequest,
		Text:         "Bad Request",
		BusinessCode: businessCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return errors.Wrap(err, "write json msg resp")
	}

	return nil
}

func NewInternalServerErr(w http.ResponseWriter) error {
	resp := HTTPError{
		Code:         http.StatusInternalServerError,
		Text:         "Internal Server error",
		BusinessCode: 0,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return errors.Wrap(err, "write json msg resp")
	}

	return nil
}

func NewNoContentResponse(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)

	if _, err := w.Write([]byte("")); err != nil {
		return errors.Wrap(err, "write empty resp")
	}

	return nil
}

func NewCreatedResponse(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write([]byte("")); err != nil {
		return errors.Wrap(err, "write empty resp")
	}

	return nil
}
