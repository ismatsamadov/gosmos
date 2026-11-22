package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"time"
)

// Service is the booking service: its data, its clock and its handlers.
type Service struct {
	Stations []Station
	Catalog  []PassView // ordered by AOS, then ID
	Guide    fs.FS      // the customer guide's pages
	Store    *MemStore  // the bookings
	Now      func() time.Time
}

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /stations", s.listStations)
	mux.HandleFunc("GET /stations/{code}", s.getStation)
	mux.Handle("GET /guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	mux.HandleFunc("GET /passes", s.listPasses)
	mux.HandleFunc("GET /passes/{sat}/{station}/{orbit}", s.getPass)
	mux.HandleFunc("POST /bookings", s.createBooking)
	mux.HandleFunc("GET /bookings/{id}", s.getBooking)
	mux.HandleFunc("DELETE /bookings/{id}", s.cancelBooking)
	return mux
}

// health answers that the service is up.
func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// writeJSON answers v as JSON with this status, or 500 in plain text when v can't be marshaled.
func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(data, '\n'))
}
