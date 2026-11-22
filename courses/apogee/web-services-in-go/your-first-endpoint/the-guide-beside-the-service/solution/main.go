package main

import "net/http"

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/stations", s.listStations)
	mux.HandleFunc("/health", s.health)
	mux.Handle("/guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	return mux
}
