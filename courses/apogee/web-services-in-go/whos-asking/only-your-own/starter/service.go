package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"time"
)

// Service is the booking service: its data, its clock and its handlers.
type Service struct {
	Stations  []Station
	Catalog   []PassView          // ordered by AOS, then ID
	Guide     fs.FS               // the customer guide's pages
	Store     Store               // the bookings
	Customers map[string]Customer // by ID
	TokenKey  []byte              // signs tokens
	Now       func() time.Time
}

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /stations", s.listStations)
	mux.HandleFunc("GET /stations/{code}", s.getStation)
	mux.Handle("GET /guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	mux.HandleFunc("POST /tokens", s.issueToken)
	mux.Handle("GET /passes", s.private(s.listPasses))
	mux.Handle("GET /passes/{sat}/{station}/{orbit}", s.private(s.getPass))
	mux.Handle("POST /bookings", s.private(s.createBooking))
	mux.Handle("GET /bookings", s.private(s.listBookings))
	mux.Handle("GET /bookings/{id}", s.private(s.getBooking))
	mux.Handle("DELETE /bookings/{id}", s.private(s.cancelBooking))
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
