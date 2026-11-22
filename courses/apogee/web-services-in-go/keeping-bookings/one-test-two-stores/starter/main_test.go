package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// The bugs Nkem planted, one in each copy of the store.
const (
	booksTwice    = "books a pass twice"
	listsUnsorted = "lists bookings in the order they were made"
	deletesNone   = "deletes nothing"
)

// mutants are the planted stores' bugs.
var mutants = []string{booksTwice, listsUnsorted, deletesNone}

// plantedStore keeps bookings as MemStore does, except for the one bug it was planted with.
type plantedStore struct {
	bug      string
	mu       sync.Mutex
	bookings []Booking // in the order they were made
	next     int
}

// Create gives b the next ID and keeps it, refusing a pass already booked unless the bug skips that check.
func (p *plantedStore) Create(ctx context.Context, b Booking) (Booking, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, old := range p.bookings {
		if p.bug != booksTwice && old.SatID == b.SatID && old.Station == b.Station && old.AOS.Equal(b.AOS) {
			return Booking{}, fmt.Errorf("%s over %s at %s: %w", b.SatID, b.Station, b.AOS.UTC().Format(time.RFC3339), ErrConflict)
		}
	}
	b.ID = fmt.Sprintf("bk-%04d", p.next)
	p.next++
	p.bookings = append(p.bookings, b)
	return b, nil
}

// Get returns the booking with this ID, or ErrNotFound.
func (p *plantedStore) Get(ctx context.Context, id string) (Booking, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, b := range p.bookings {
		if b.ID == id {
			return b, nil
		}
	}
	return Booking{}, fmt.Errorf("booking %s: %w", id, ErrNotFound)
}

// List returns a copy of the bookings by AOS and then ID, unless the bug leaves them unsorted.
func (p *plantedStore) List(ctx context.Context) ([]Booking, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	list := slices.Clone(p.bookings)
	if p.bug != listsUnsorted {
		slices.SortFunc(list, func(a, b Booking) int {
			if c := a.AOS.Compare(b.AOS); c != 0 {
				return c
			}
			return strings.Compare(a.ID, b.ID)
		})
	}
	return list, nil
}

// Delete removes the booking with this ID, or returns ErrNotFound; with the bug it only says it removed it.
func (p *plantedStore) Delete(ctx context.Context, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	i := slices.IndexFunc(p.bookings, func(b Booking) bool { return b.ID == id })
	if i < 0 {
		return fmt.Errorf("booking %s: %w", id, ErrNotFound)
	}
	if p.bug != deletesNone {
		p.bookings = slices.Delete(p.bookings, i, i+1)
	}
	return nil
}

func TestMemStoreKeepsTheContract(t *testing.T) {
	checkStore(t, func() Store { return NewMemStore() })
}

// TestMutant runs checkStore against the planted store APOGEE_MUTANT names. In a normal run it
// does nothing: TestChecksCatchMutants runs it in a fresh copy of this test binary.
func TestMutant(t *testing.T) {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return
	}
	if !slices.Contains(mutants, bug) {
		t.Fatalf("no planted store that %q", bug)
	}
	checkStore(t, func() Store { return &plantedStore{bug: bug, next: 4501} })
}

// TestChecksCatchMutants runs each planted store in a fresh copy of this test binary. The run
// must fail, and the failure must belong to a subtest, so that the log names the case.
func TestChecksCatchMutants(t *testing.T) {
	for _, bug := range mutants {
		cmd := exec.Command(os.Args[0], "-test.run=^TestMutant$", "-test.timeout=5s")
		cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+bug)
		out, err := cmd.CombinedOutput()
		var exit *exec.ExitError
		switch {
		case err == nil:
			t.Errorf("the store that %s survived: checkStore reported nothing against it", bug)
		case !errors.As(err, &exit):
			t.Fatalf("running the store that %s: %v", bug, err)
		case !strings.Contains(string(out), "--- FAIL: TestMutant/"):
			t.Errorf("the store that %s failed no subtest, only TestMutant; want the failing case's subtest to fail", bug)
		}
	}
}
