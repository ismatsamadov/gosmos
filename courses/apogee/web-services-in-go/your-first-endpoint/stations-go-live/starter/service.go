package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"sync"
)

// Service is the booking service: its data and its handlers.
type Service struct {
	Stations []Station
	Guide    fs.FS // the customer guide's pages
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

// health answers that the service is up.
func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// routes registers the service's handlers on a new mux.
func (s *Service) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/stations", s.listStations)
	mux.HandleFunc("/health", s.health)
	mux.Handle("/guide/", http.StripPrefix("/guide/", http.FileServerFS(s.Guide)))
	return mux
}

// stopFunc returns a function that closes srv and waits for its Serve, which reports on served, to return.
func stopFunc(srv *http.Server, served <-chan error) func() error {
	var once sync.Once
	var err error
	return func() error {
		once.Do(func() {
			err = srv.Close()
			if serveErr := <-served; !errors.Is(serveErr, http.ErrServerClosed) {
				err = serveErr
			}
		})
		return err
	}
}
