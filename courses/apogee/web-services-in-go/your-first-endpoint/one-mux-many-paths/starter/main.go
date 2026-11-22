package main

import "net/http"

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/stations", s.listStations)
	return mux
}

// health answers that the service is up.
func (s *Service) health(w http.ResponseWriter, r *http.Request) {
}
