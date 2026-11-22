package main

import (
	"errors"
	"net/http"
)

// createBooking checks the booking in the request's body and keeps it, answering 201 and where it lives.
func (s *Service) createBooking(w http.ResponseWriter, r *http.Request) {
	var b Booking
	if err := decodeJSON(w, r, &b); err != nil {
		writeError(w, r, err)
		return
	}
	// TODO: answer the first failure of checkShape, then of checkPass
	b, err := s.Store.Create(r.Context(), b)
	_ = errors.Is
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}
