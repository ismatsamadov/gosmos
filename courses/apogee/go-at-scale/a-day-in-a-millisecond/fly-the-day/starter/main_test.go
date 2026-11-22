package main

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

// takeOverDay is midnight UTC on Thursday 28 June 2029, where the tests set a bubble's clock.
var takeOverDay = time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)

// clockAt returns a clock that reads day when the bubble starts and moves with its fake time.
func clockAt(day time.Time) func() time.Time {
	shift := day.Sub(time.Now())
	return func() time.Time { return time.Now().Add(shift).UTC() }
}

// takeOverDraft, 28 June as passplan 2 plans it from scratch, is in a file the lesson keeps hidden.

func TestFlyDay(t *testing.T) {
	day := takeOverDraft()
	if len(day) != 2020 {
		t.Fatalf("takeOverDraft() has %d contacts, want 2020", len(day))
	}
	noon := takeOverDay.Add(12 * time.Hour)
	tests := []struct {
		name    string
		cancel  bool // cancel the context at 12:00
		wantErr error
		flown   int
	}{
		{"take-over day's 2,020 contacts", false, nil, 2020},
		{"take-over day, canceled at 12:00", true, context.Canceled, 992},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				now := clockAt(takeOverDay)
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				if tt.cancel {
					time.AfterFunc(noon.Sub(now()), cancel)
				}
				var flown []Contact
				var at []time.Time // when each was flown
				err := flyDay(ctx, now, day, func(c Contact) {
					flown = append(flown, c)
					at = append(at, now())
				})
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("flyDay(ctx, now, day, fly) = %v, want %v", err, tt.wantErr)
				}
				if tt.cancel && !now().Equal(noon) {
					t.Errorf("flyDay(ctx, now, day, fly) returned at %s, want 12:00:00, when its context ended", now().Format(time.TimeOnly))
				}
				if len(flown) != tt.flown {
					t.Errorf("flyDay(ctx, now, day, fly) flew %d contacts, want %d", len(flown), tt.flown)
				}
				for i, c := range flown[:min(len(flown), len(day))] {
					if c.ID != day[i].ID {
						t.Errorf("flyDay(ctx, now, day, fly) contact %d = %s, want %s, in the day's order", i+1, c.ID, day[i].ID)
						break
					}
					if !at[i].Equal(c.AOS) {
						t.Errorf("flyDay(ctx, now, day, fly) flew %s at %s, want %s, its AOS", c.ID, at[i].Format(time.TimeOnly), c.AOS.Format(time.TimeOnly))
						break
					}
				}
			})
		})
	}
}
