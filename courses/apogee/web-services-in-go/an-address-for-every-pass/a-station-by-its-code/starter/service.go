package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
)

// Service is the booking service: its data and its handlers.
type Service struct {
	Stations []Station
	Guide    fs.FS // the customer guide's pages
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
