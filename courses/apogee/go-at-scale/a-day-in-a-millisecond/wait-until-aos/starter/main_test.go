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

func TestWaitUntil(t *testing.T) {
	aos := takeOverDay.Add(4*time.Hour + 12*time.Minute + 25*time.Second)
	tests := []struct {
		name    string
		at      time.Time     // the instant to wait for
		limit   time.Duration // the context's deadline, after the start; 0 for none
		wantErr error
		wantAt  time.Time
	}{
		{"Furrow-S over Corith Point at 04:12:25", aos, 0, nil, aos},
		{"an AOS already past", takeOverDay.Add(-100 * time.Second), 0, nil, takeOverDay},
		{"a 30 s deadline on a wait of an hour", takeOverDay.Add(time.Hour), 30 * time.Second, context.DeadlineExceeded, takeOverDay.Add(30 * time.Second)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				now := clockAt(takeOverDay)
				ctx := t.Context()
				if tt.limit > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tt.limit)
					defer cancel()
				}
				err := waitUntil(ctx, now, tt.at)
				if got := now(); !errors.Is(err, tt.wantErr) || !got.Equal(tt.wantAt) {
					t.Errorf("waitUntil(ctx, now, %s) = %v at %s, want %v at %s", tt.at.Format(time.TimeOnly),
						err, got.Format(time.TimeOnly), tt.wantErr, tt.wantAt.Format(time.TimeOnly))
				}
			})
		})
	}
}
