package main

import (
	"encoding/json"
	"net/http"
)

// createBooking keeps the booking in the request's body and answers 201 with it, its new ID included.
func (s *Service) createBooking(w http.ResponseWriter, r *http.Request) {
	var b Booking
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	b, err := s.Store.Create(r.Context(), b)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}
