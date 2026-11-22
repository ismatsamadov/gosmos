package main

import (
	"encoding/json"
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

// createBooking keeps the booking in the request's body and answers 201 with it, its new ID included.
func (s *Service) createBooking(w http.ResponseWriter, r *http.Request) {
	var b Booking
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeError(w, r, &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: "request body: " + err.Error()})
		return
	}
	b, err := s.Store.Create(r.Context(), b)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}
