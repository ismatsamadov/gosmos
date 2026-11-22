package main

import (
	"fmt"
	"net/http"
)

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /stations", s.listStations)
	mux.Handle("GET /guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	return mux
}

// getStation answers with the station whose code is in the path, or 404.
func (s *Service) getStation(w http.ResponseWriter, r *http.Request) {
	code := "" // TODO: read the code from the path
	st, ok := findStation(s.Stations, code)
	if !ok {
		http.Error(w, fmt.Sprintf("no station %q", code), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, st)
}
