package main

import (
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

// getBooking answers the booking whose ID the path names.
func (s *Service) getBooking(w http.ResponseWriter, r *http.Request) {
	b, err := s.Store.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

// cancelBooking removes the booking whose ID the path names, answering 204 with no body.
func (s *Service) cancelBooking(w http.ResponseWriter, r *http.Request) {
	if err := s.Store.Delete(r.Context(), r.PathValue("id")); err != nil {
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
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Location", "/bookings/"+b.ID)
	writeJSON(w, http.StatusCreated, b)
}

// listBookings answers every booking, ordered by AOS and then ID.
func (s *Service) listBookings(w http.ResponseWriter, r *http.Request) {
	list, err := s.Store.List(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
