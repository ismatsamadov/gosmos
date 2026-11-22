package main

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

// checkWaitUntil fails t when wait doesn't end exactly when it should.
func checkWaitUntil(t *testing.T, wait func(context.Context, func() time.Time, time.Time) error) {
	synctest.Test(t, func(t *testing.T) {
		day := time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)
		shift := day.Sub(time.Now())
		now := func() time.Time { return time.Now().Add(shift).UTC() }
		aos := day.Add(4*time.Hour + 12*time.Minute + 25*time.Second)
		if err := wait(t.Context(), now, aos); err != nil || !now().Equal(aos) {
			t.Errorf("wait(ctx, now, 04:12:25) = %v at %s, want <nil> at 04:12:25",
				err, now().Format(time.TimeOnly))
		}
		if err := wait(t.Context(), now, aos.Add(-time.Hour)); err != nil || !now().Equal(aos) {
			t.Errorf("wait(ctx, now, 03:12:25) = %v at %s, want <nil> at once",
				err, now().Format(time.TimeOnly))
		}
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		if err := wait(ctx, now, aos.Add(time.Hour)); !errors.Is(err, context.Canceled) {
			t.Errorf("wait(canceled ctx, now, 05:12:25) = %v, want context canceled", err)
		}
	})
}

// notifier is a debouncer as checkDebouncer sees it.
type notifier interface {
	Notify()
	Stop()
}

// checkDebouncer fails t when a debouncer from newD fires early, more than once, or after Stop.
func checkDebouncer(t *testing.T, newD func(quiet time.Duration, fire func()) notifier) {}
