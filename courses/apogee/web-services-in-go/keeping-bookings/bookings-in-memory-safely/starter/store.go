package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

// The errors a Store wraps, with what the booking or the pass was.
var (
	ErrNotFound = errors.New("not found")      // wrapped: "booking bk-5499: not found"
	ErrConflict = errors.New("already booked") // wrapped: "furrow-a over SQT at 2027-12-06T14:02:10Z: already booked"
)

// Store keeps bookings. Its methods are safe to call from many goroutines at once.
type Store interface {
	// Create gives b the next ID and keeps it. It returns ErrConflict when a booking
	// of the same pass (satellite, station and AOS) is already kept.
	Create(ctx context.Context, b Booking) (Booking, error)
	// Get returns the booking with this ID, or ErrNotFound.
	Get(ctx context.Context, id string) (Booking, error)
	// List returns every booking, ordered by AOS and then ID, in a new slice.
	List(ctx context.Context) ([]Booking, error)
	// Delete removes the booking with this ID, or returns ErrNotFound.
	Delete(ctx context.Context, id string) error
}

// firstNumber is the number of the service's first booking, bk-4501.
const firstNumber = 4501

// MemStore keeps bookings in a map behind a mutex.
type MemStore struct {
	mu       sync.Mutex
	bookings map[string]Booking // by ID
	next     int                // the number of the next ID
}

// NewMemStore returns an empty store whose first booking gets bk-4501.
func NewMemStore() *MemStore {
	return &MemStore{bookings: map[string]Booking{}, next: firstNumber}
}

// Create gives b the next ID and keeps it, unless a booking of the same pass is already kept.
func (m *MemStore) Create(ctx context.Context, b Booking) (Booking, error) {
	_ = fmt.Errorf
	_ = time.RFC3339
	return Booking{}, nil
}

// Get returns the booking with this ID, or ErrNotFound.
func (m *MemStore) Get(ctx context.Context, id string) (Booking, error) {
	return Booking{}, nil
}

// List returns every booking, ordered by AOS and then ID, in a new slice.
func (m *MemStore) List(ctx context.Context) ([]Booking, error) {
	_ = slices.SortFunc[[]Booking, Booking]
	_ = strings.Compare
	return nil, nil
}

// Delete removes the booking with this ID, or returns ErrNotFound.
func (m *MemStore) Delete(ctx context.Context, id string) error {
	return nil
}
