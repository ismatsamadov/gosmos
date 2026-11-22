package main

import "net/http"

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /stations", s.listStations)
	mux.HandleFunc("GET /stations/{code}", s.getStation)
	mux.Handle("GET /guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	mux.HandleFunc("GET /passes", s.listPasses)
	mux.HandleFunc("GET /passes/{sat}/{station}/{orbit}", s.getPass)
	mux.HandleFunc("GET /bookings/{id}", s.getBooking)
	mux.HandleFunc("DELETE /bookings/{id}", s.cancelBooking)
	return mux
}
