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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

// cancelBooking removes the booking whose ID is in the path and answers 204, or 404 when there is none.
func (s *Service) cancelBooking(w http.ResponseWriter, r *http.Request) {
	err := s.Store.Delete(r.Context(), r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
