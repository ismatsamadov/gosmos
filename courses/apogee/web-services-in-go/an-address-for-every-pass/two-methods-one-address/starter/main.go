package main

import (
	"errors"
	"net/http"
)

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /stations", s.listStations)
	mux.HandleFunc("GET /stations/{code}", s.getStation)
	mux.HandleFunc("GET /passes", s.listPasses)
	mux.HandleFunc("GET /passes/{sat}/{station}/{orbit}", s.getPass)
	mux.HandleFunc("GET /bookings/{id}", s.getBooking)
	mux.Handle("GET /guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	return mux
}

// cancelBooking removes the booking whose ID is in the path and answers 204, or 404 when there is none.
func (s *Service) cancelBooking(w http.ResponseWriter, r *http.Request) {
	_ = errors.Is
	// TODO: delete the booking from s.Store, then answer 204 No Content
}
