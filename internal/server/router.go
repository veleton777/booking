//nolint:godot
package server

import (
	"github.com/go-chi/chi/v5"
)

// @title Swagger Currency API
// @version 1.0
// @description This is a server Currency API

// @host localhost:8080
// @BasePath /api
func (s *API) routes(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/booking", s.bookingServer.Booking)
	})
}
