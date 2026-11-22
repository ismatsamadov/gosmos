package main

import (
	"encoding/json"
	"net/http"
)

// createBooking keeps the booking in the request's body and answers 201 with it, its new ID included.
func (s *Service) createBooking(w http.ResponseWriter, r *http.Request) {
	var b Booking
	_ = json.NewDecoder
	// TODO: decode r.Body into b, then keep b with s.Store.Create
	writeJSON(w, http.StatusCreated, b)
}
