package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
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
	if err := ctx.Err(); err != nil {
		return Booking{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, old := range m.bookings {
		if old.SatID == b.SatID && old.Station == b.Station && old.AOS.Equal(b.AOS) {
			return Booking{}, fmt.Errorf("%s over %s at %s: %w", b.SatID, b.Station, b.AOS.UTC().Format(time.RFC3339), ErrConflict)
		}
	}
	b.ID = fmt.Sprintf("bk-%04d", m.next)
	m.next++
	m.bookings[b.ID] = b
	return b, nil
}

// Get returns the booking with this ID, or ErrNotFound.
func (m *MemStore) Get(ctx context.Context, id string) (Booking, error) {
	if err := ctx.Err(); err != nil {
		return Booking{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bookings[id]
	if !ok {
		return Booking{}, fmt.Errorf("booking %s: %w", id, ErrNotFound)
	}
	return b, nil
}

// List returns every booking, ordered by AOS and then ID, in a new slice.
func (m *MemStore) List(ctx context.Context) ([]Booking, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	list := make([]Booking, 0, len(m.bookings))
	for _, b := range m.bookings {
		list = append(list, b)
	}
	slices.SortFunc(list, func(a, b Booking) int {
		if c := a.AOS.Compare(b.AOS); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	return list, nil
}

// Delete removes the booking with this ID, or returns ErrNotFound.
func (m *MemStore) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bookings[id]; !ok {
		return fmt.Errorf("booking %s: %w", id, ErrNotFound)
	}
	delete(m.bookings, id)
	return nil
}

// load keeps saved bookings under their own IDs and moves the next ID past the highest of them.
func (m *MemStore) load(saved []Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, b := range saved {
		digits, ok := strings.CutPrefix(b.ID, "bk-")
		n, err := strconv.Atoi(digits)
		if !ok || err != nil {
			return fmt.Errorf("booking ID %q is not bk- and a number", b.ID)
		}
		if _, dup := m.bookings[b.ID]; dup {
			return fmt.Errorf("booking %s is kept twice", b.ID)
		}
		m.bookings[b.ID] = b
		m.next = max(m.next, n+1)
	}
	return nil
}
