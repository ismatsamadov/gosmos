package main

import (
	"fmt"
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

// getBooking answers the caller's booking whose ID the path names.
func (s *Service) getBooking(w http.ResponseWriter, r *http.Request) {
	b, err := s.findBooking(r)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

// cancelBooking removes the caller's booking whose ID the path names, answering 204 with no body.
func (s *Service) cancelBooking(w http.ResponseWriter, r *http.Request) {
	b, err := s.findBooking(r)
	if err == nil {
		err = s.Store.Delete(r.Context(), b.ID)
	}
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// findBooking returns the caller's booking whose ID the path names; another customer's is not found, exactly like a missing one.
func (s *Service) findBooking(r *http.Request) (Booking, error) {
	c, _ := customerFrom(r.Context())
	id := r.PathValue("id")
	b, err := s.Store.Get(r.Context(), id)
	if err == nil && !owns(c, b.SatID) {
		return Booking{}, fmt.Errorf("booking %s: %w", id, ErrNotFound)
	}
	return b, err
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
	if c, _ := customerFrom(r.Context()); !owns(c, b.SatID) {
		writeError(w, r, &APIError{Status: http.StatusForbidden, Code: "forbidden", Message: b.SatID + " is not one of your satellites"})
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

// listBookings answers the caller's bookings, ordered by AOS and then ID.
func (s *Service) listBookings(w http.ResponseWriter, r *http.Request) {
	list, err := s.Store.List(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	c, _ := customerFrom(r.Context())
	mine := []Booking{}
	for _, b := range list {
		if owns(c, b.SatID) {
			mine = append(mine, b)
		}
	}
	writeJSON(w, http.StatusOK, mine)
}
