package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"sync/atomic"
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
	Limiter   *Limiter
	Timeouts  Timeouts
	Keys      *keyStore // idempotency keys
	Now       func() time.Time
	draining  atomic.Bool // set while the service shuts down
}

// Config is everything the service is built from.
type Config struct {
	Dir       string // the data folder, which holds bookings.json and outbox.json
	Stations  []Station
	Catalog   []PassView
	Guide     fs.FS // embedded in production, an fstest.MapFS in tests
	Customers []Customer
	TokenKey  []byte
	Now       func() time.Time
	Rand      *rand.Rand // jitter, seeded in tests
	Log       io.Writer  // the access log
	Timeouts  Timeouts
}

// newService builds the service from cfg, keeping its bookings in cfg.Dir.
func newService(cfg Config) (*Service, error) {
	store, err := OpenFileStore(filepath.Join(cfg.Dir, "bookings.json"))
	if err != nil {
		return nil, err
	}
	customers := make(map[string]Customer, len(cfg.Customers))
	for _, c := range cfg.Customers {
		customers[c.ID] = c
	}
	return &Service{
		Stations:  cfg.Stations,
		Catalog:   cfg.Catalog,
		Guide:     cfg.Guide,
		Store:     store,
		Customers: customers,
		TokenKey:  cfg.TokenKey,
		Limiter:   newLimiter(cfg.Now),
		Timeouts:  cfg.Timeouts,
		Keys:      newKeyStore(cfg.Now),
		Now:       cfg.Now,
	}, nil
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
	book := idempotent(s.Keys)(http.HandlerFunc(s.createBooking))
	mux.Handle("POST /bookings", s.private(book.ServeHTTP))
	mux.Handle("GET /bookings", s.private(s.listBookings))
	mux.Handle("GET /bookings/{id}", s.private(s.getBooking))
	mux.Handle("DELETE /bookings/{id}", s.private(s.cancelBooking))
	return mux
}

// health answers GET /health: ok, or 503 draining while the service shuts down.
func (s *Service) health(w http.ResponseWriter, r *http.Request) {
	if s.draining.Load() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "draining"})
		return
	}
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
