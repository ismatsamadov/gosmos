package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// jan3 returns a UTC time on Monday 3 January 2028.
func jan3(hour, minute, second int) time.Time {
	return time.Date(2028, 1, 3, hour, minute, second, 0, time.UTC)
}

// booking returns a made-up booking of sat over station that rises at aos and lasts ten minutes:
// the store keeps any booking, so its rules can be tried on passes no catalog lists.
func booking(sat, station string, aos time.Time) Booking {
	return Booking{SatID: sat, Station: station, AOS: aos, LOS: aos.Add(10 * time.Minute)}
}

// Furrowlight's predicted passes over Salar Quieto on 3 January.
var (
	furrowA = Booking{SatID: "furrow-a", Station: "SQT", AOS: jan3(14, 11, 55), LOS: jan3(14, 21, 59)}
	furrowB = Booking{SatID: "furrow-b", Station: "SQT", AOS: jan3(14, 59, 32), LOS: jan3(15, 9, 5)}
)

// Two made-up bookings that rise at the same second as Furrow-A's pass: one of another
// satellite at the same station, one at another station.
var (
	brine3 = booking("brine-3", "SQT", furrowA.AOS)
	pipit1 = booking("pipit-1", "PEL", furrowA.AOS)
)

func TestCreate(t *testing.T) {
	m := NewMemStore()
	for i, b := range []Booking{furrowA, furrowB, pipit1} {
		want := fmt.Sprintf("bk-%d", 4501+i)
		got, err := m.Create(context.Background(), b)
		if err != nil || got.ID != want {
			t.Errorf("Create(%s over %s) = %q, %v, want %q, nil", b.SatID, b.Station, got.ID, err, want)
			continue
		}
		if got.SatID != b.SatID || got.Station != b.Station || !got.AOS.Equal(b.AOS) || !got.LOS.Equal(b.LOS) {
			t.Errorf("Create(%s over %s) = %+v, want the booking's own fields with ID %q", b.SatID, b.Station, got, want)
		}
	}
}

func TestCreateConflict(t *testing.T) {
	ctx := context.Background()
	m := NewMemStore()
	if _, err := m.Create(ctx, furrowA); err != nil {
		t.Fatalf("Create(furrow-a over SQT) error = %v, want nil", err)
	}
	again := furrowA
	again.Note = "SECOND TRY, SAME PASS"
	_, err := m.Create(ctx, again)
	const want = "furrow-a over SQT at 2028-01-03T14:11:55Z: already booked"
	if !errors.Is(err, ErrConflict) || err.Error() != want {
		t.Errorf("Create(the same pass again) error = %v, want %q, wrapping ErrConflict", err, want)
	}
	if _, err := m.Create(ctx, brine3); err != nil {
		t.Errorf("Create(brine-3 over SQT at the same second) error = %v, want nil: it is another satellite's pass", err)
	}
}

func TestGetAndDelete(t *testing.T) {
	ctx := context.Background()
	m := NewMemStore()
	made, err := m.Create(ctx, furrowA)
	if err != nil {
		t.Fatalf("Create(furrow-a over SQT) error = %v, want nil", err)
	}
	if got, err := m.Get(ctx, made.ID); err != nil || got != made {
		t.Errorf("Get(%q) = %+v, %v, want %+v, nil", made.ID, got, err, made)
	}
	if err := m.Delete(ctx, made.ID); err != nil {
		t.Errorf("Delete(%q) = %v, want nil", made.ID, err)
	}
	if _, err := m.Get(ctx, made.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(%q) after Delete error = %v, want ErrNotFound", made.ID, err)
	}
	const want = "booking bk-5499: not found"
	if _, err := m.Get(ctx, "bk-5499"); !errors.Is(err, ErrNotFound) || err.Error() != want {
		t.Errorf("Get(%q) error = %v, want %q, wrapping ErrNotFound", "bk-5499", err, want)
	}
	if err := m.Delete(ctx, "bk-5499"); !errors.Is(err, ErrNotFound) || err.Error() != want {
		t.Errorf("Delete(%q) error = %v, want %q, wrapping ErrNotFound", "bk-5499", err, want)
	}
}

func TestList(t *testing.T) {
	ctx := context.Background()
	m := NewMemStore()
	for _, b := range []Booking{furrowB, furrowA, brine3, pipit1} { // bk-4501 to bk-4504
		if _, err := m.Create(ctx, b); err != nil {
			t.Fatalf("Create(%s over %s) error = %v, want nil", b.SatID, b.Station, err)
		}
	}
	want := []string{"bk-4502", "bk-4503", "bk-4504", "bk-4501"} // by AOS, then ID
	list, err := m.List(ctx)
	if err != nil || len(list) != len(want) {
		t.Fatalf("List() = %d bookings, %v, want %d, nil", len(list), err, len(want))
	}
	for i := range want {
		if list[i].ID != want[i] {
			t.Errorf("List()[%d].ID = %q, want %q (AOS first, then ID)", i, list[i].ID, want[i])
		}
	}
	list[0].Note = "CHANGED BY A CALLER"
	again, err := m.List(ctx)
	if err != nil || len(again) != len(want) {
		t.Fatalf("List() again = %d bookings, %v, want %d, nil", len(again), err, len(want))
	}
	if again[0].Note != "" {
		t.Errorf("after a caller changed List's result, List()[0].Note = %q, want %q: the store must hand out copies", again[0].Note, "")
	}
}

func TestManyBookingsAtOnce(t *testing.T) {
	ctx := context.Background()
	t.Run("50 customers book 50 passes", func(t *testing.T) {
		m := NewMemStore()
		ids := make([]string, 50)
		var wg sync.WaitGroup
		for i := range ids { // 50 made-up bookings, each rising 95 minutes after the one before
			wg.Go(func() {
				b, _ := m.Create(ctx, booking("pipit-1", "PEL", jan3(0, 0, 0).Add(time.Duration(i)*95*time.Minute)))
				ids[i] = b.ID
			})
		}
		wg.Wait()
		seen := map[string]bool{}
		for _, id := range ids {
			if id != "" {
				seen[id] = true
			}
		}
		if len(seen) != len(ids) {
			t.Errorf("50 bookings of 50 passes made at once got %d different IDs, want %d", len(seen), len(ids))
		}
	})
	t.Run("50 customers book one pass", func(t *testing.T) {
		m := NewMemStore()
		errs := make([]error, 50)
		var wg sync.WaitGroup
		for i := range errs {
			wg.Go(func() {
				_, errs[i] = m.Create(ctx, furrowA)
			})
		}
		wg.Wait()
		won := 0
		for _, err := range errs {
			switch {
			case err == nil:
				won++
			case !errors.Is(err, ErrConflict):
				t.Errorf("Create(the pass everyone wants) error = %v, want nil or ErrConflict", err)
			}
		}
		if won != 1 {
			t.Errorf("50 bookings of one pass at once: %d succeeded, want exactly 1", won)
		}
	})
}

func TestCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	m := NewMemStore()
	if _, err := m.Create(ctx, furrowA); !errors.Is(err, context.Canceled) {
		t.Errorf("Create(canceled context) error = %v, want context.Canceled", err)
	}
	if _, err := m.Get(ctx, "bk-4501"); !errors.Is(err, context.Canceled) {
		t.Errorf("Get(canceled context) error = %v, want context.Canceled", err)
	}
	if _, err := m.List(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("List(canceled context) error = %v, want context.Canceled", err)
	}
	if err := m.Delete(ctx, "bk-4501"); !errors.Is(err, context.Canceled) {
		t.Errorf("Delete(canceled context) error = %v, want context.Canceled", err)
	}
}
