package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

// checkStore checks that every store open returns keeps the promises of the Store interface.
func checkStore(t *testing.T, open func() Store) {
	ctx := context.Background()
	early := Booking{SatID: "furrow-a", Station: "SQT",
		AOS: time.Date(2028, 1, 3, 14, 11, 55, 0, time.UTC), LOS: time.Date(2028, 1, 3, 14, 21, 59, 0, time.UTC)}
	late := Booking{SatID: "furrow-b", Station: "SQT",
		AOS: time.Date(2028, 1, 3, 14, 59, 32, 0, time.UTC), LOS: time.Date(2028, 1, 3, 15, 9, 5, 0, time.UTC)}
	create := func(t *testing.T, s Store, b Booking) Booking {
		t.Helper()
		made, err := s.Create(ctx, b)
		if err != nil {
			t.Fatalf("Create(%s at %s) error = %v, want nil", b.SatID, b.AOS.Format(time.RFC3339), err)
		}
		return made
	}

	t.Run("IDs count up from bk-4501", func(t *testing.T) {
		s := open()
		if got := create(t, s, late).ID; got != "bk-4501" {
			t.Errorf("first Create ID = %q, want %q", got, "bk-4501")
		}
		if got := create(t, s, early).ID; got != "bk-4502" {
			t.Errorf("second Create ID = %q, want %q", got, "bk-4502")
		}
	})
	t.Run("Get returns what Create kept", func(t *testing.T) {
		s := open()
		made := create(t, s, early)
		if got, err := s.Get(ctx, made.ID); err != nil || got != made {
			t.Errorf("Get(%q) = %v, %v, want %v, nil", made.ID, got, err, made)
		}
	})
	t.Run("a pass is booked only once", func(t *testing.T) {
		s := open()
		create(t, s, early)
		if _, err := s.Create(ctx, early); !errors.Is(err, ErrConflict) {
			t.Errorf("Create(the same pass again) error = %v, want ErrConflict", err)
		}
	})
	t.Run("Get after Delete finds nothing", func(t *testing.T) {
		s := open()
		made := create(t, s, early)
		if err := s.Delete(ctx, made.ID); err != nil {
			t.Fatalf("Delete(%q) = %v, want nil", made.ID, err)
		}
		if _, err := s.Get(ctx, made.ID); !errors.Is(err, ErrNotFound) {
			t.Errorf("Get(%q) after Delete error = %v, want ErrNotFound", made.ID, err)
		}
	})
	t.Run("List orders by AOS", func(t *testing.T) {
		s := open()
		second, first := create(t, s, late), create(t, s, early)
		list, err := s.List(ctx)
		if err != nil || len(list) != 2 {
			t.Fatalf("List() = %d bookings, %v, want 2, nil", len(list), err)
		}
		if list[0].ID != first.ID || list[1].ID != second.ID {
			t.Errorf("List() IDs = %s, %s, want %s, %s", list[0].ID, list[1].ID, first.ID, second.ID)
		}
	})
}
