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
	if err := checkShape(b); err != nil {
		writeError(w, r, err)
		return
	}
	if err := checkPass(b, s.Catalog, s.Now()); err != nil {
		writeError(w, r, err)
		return
	}
	b, err := s.Store.Create(r.Context(), b)
	if errors.Is(err, ErrConflict) {
		err = &APIError{Status: http.StatusConflict, Code: "conflict", Message: err.Error()}
	}
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Location", "/bookings/"+b.ID)
	writeJSON(w, http.StatusCreated, b)
}
