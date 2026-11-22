package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"
)

// mutants are Nkem's broken copies of waitUntil, each with one bug planted.
var mutants = []struct {
	bug  string
	wait func(context.Context, func() time.Time, time.Time) error
}{
	{"wakes a second early", func(ctx context.Context, now func() time.Time, at time.Time) error {
		return waitUntil(ctx, now, at.Add(-time.Second))
	}},
	{"ignores its context", func(_ context.Context, now func() time.Time, at time.Time) error {
		time.Sleep(at.Sub(now()))
		return nil
	}},
	{"returns at once", func(ctx context.Context, _ func() time.Time, _ time.Time) error {
		return ctx.Err()
	}},
	{"waits a minute for an instant already past", func(ctx context.Context, now func() time.Time, at time.Time) error {
		if !now().Before(at) {
			time.Sleep(time.Minute)
		}
		return waitUntil(ctx, now, at)
	}},
}

// waiter returns the wait this run checks: your waitUntil, or the mutant APOGEE_MUTANT names.
func waiter(t *testing.T) func(context.Context, func() time.Time, time.Time) error {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return waitUntil
	}
	for _, m := range mutants {
		if m.bug == bug {
			return m.wait
		}
	}
	t.Fatalf("no mutant %q", bug)
	return nil
}

func TestCheckWaitUntil(t *testing.T) {
	wait := waiter(t)
	start := time.Now()       // the real clock; a bubble's fake one starts in 2000
	var realClock atomic.Bool // set when wait runs on the real clock, outside a bubble
	// Twenty runs: a check that passes only by a select's choice fails some of them.
	for range 20 {
		checkWaitUntil(t, func(ctx context.Context, now func() time.Time, at time.Time) error {
			if d := time.Since(start); d >= 0 && d < time.Minute {
				realClock.Store(true)
			}
			return wait(ctx, now, at)
		})
		if t.Failed() || realClock.Load() {
			break
		}
	}
	if realClock.Load() && os.Getenv("APOGEE_MUTANT") == "" {
		t.Errorf("checkWaitUntil called wait on the real clock, want every call inside synctest.Test")
	}
}

func TestChecksCatchMutants(t *testing.T) {
	for _, m := range mutants {
		t.Run("the mutant that "+m.bug, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestCheckWaitUntil$", "-test.timeout=5s")
			cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
			err := cmd.Run()
			var exit *exec.ExitError
			switch {
			case err == nil:
				t.Errorf("the mutant that %s survived: checkWaitUntil reported nothing", m.bug)
			case !errors.As(err, &exit):
				t.Fatalf("running the mutant that %s: %v", m.bug, err)
			}
		})
	}
}
