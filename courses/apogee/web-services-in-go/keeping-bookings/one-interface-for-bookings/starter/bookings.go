package main

import (
	"errors"
	"net/http"
	"time"
)

// Booking is one reserved pass.
type Booking struct {
	ID      string    `json:"id"` // "bk-4501", given by the service
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
	LOS     time.Time `json:"los"`
	Note    string    `json:"note,omitempty"`
}

// getBooking answers with the booking whose ID is in the path, or 404.
func (s *Service) getBooking(w http.ResponseWriter, r *http.Request) {
	b, err := s.Store.Get(r.Context(), r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		err = &APIError{Status: http.StatusNotFound, Code: "not_found", Message: err.Error()}
	}
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

// cancelBooking removes the booking whose ID is in the path and answers 204, or 404 when there is none.
func (s *Service) cancelBooking(w http.ResponseWriter, r *http.Request) {
	err := s.Store.Delete(r.Context(), r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		err = &APIError{Status: http.StatusNotFound, Code: "not_found", Message: err.Error()}
	}
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
