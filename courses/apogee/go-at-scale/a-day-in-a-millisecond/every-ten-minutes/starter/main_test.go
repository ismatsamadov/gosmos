package main

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"
	"time"
)

// clock writes d, a time after the bubble's start, as hours, minutes and seconds: 24:00:00 for a day.
func clock(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
}

func TestEvery(t *testing.T) {
	tests := []struct {
		name  string
		d     time.Duration // every's period
		run   time.Duration // how long each call of f takes
		stop  time.Duration // when the test cancels the context
		calls int
		ended time.Duration // when every should return
	}{
		{"a day of replans, ten minutes apart", 10 * time.Minute, 0, 24*time.Hour + time.Second, 144, 24*time.Hour + time.Second},
		{"a day of replans that take 20 seconds each", 10 * time.Minute, 20 * time.Second, 24*time.Hour + time.Second, 144, 24*time.Hour + 20*time.Second},
		{"a day canceled at 09:47", 10 * time.Minute, 0, 9*time.Hour + 47*time.Minute, 58, 9*time.Hour + 47*time.Minute},
		{"a day of replans six hours apart", 6 * time.Hour, 0, 24*time.Hour + time.Second, 4, 24*time.Hour + time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				start := time.Now()
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				time.AfterFunc(tt.stop, cancel)
				var calls []time.Duration
				every(ctx, tt.d, func(at time.Time) {
					if ctx.Err() != nil {
						t.Fatalf("every(ctx, %v, f) called f at %s, after its context ended at %s", tt.d, clock(time.Since(start)), clock(tt.stop))
					}
					calls = append(calls, at.Sub(start))
					time.Sleep(tt.run)
				})
				if len(calls) != tt.calls {
					t.Errorf("every(ctx, %v, f): calls = %d, want %d", tt.d, len(calls), tt.calls)
				}
				for i, got := range calls {
					if want := time.Duration(i+1) * tt.d; got != want {
						t.Errorf("every(ctx, %v, f) call %d at %s, want %s", tt.d, i+1, clock(got), clock(want))
						break
					}
				}
				if ended := time.Since(start); ended != tt.ended {
					t.Errorf("every(ctx, %v, f) returned at %s, want %s", tt.d, clock(ended), clock(tt.ended))
				}
			})
		})
	}
}
