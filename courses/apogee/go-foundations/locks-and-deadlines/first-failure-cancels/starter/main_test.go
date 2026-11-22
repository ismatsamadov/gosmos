package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

// check returns one station's pre-pass check. With fail set, it fails at once; otherwise it
// waits until ctx is canceled. Either way, it reports its station on finished as it returns.
func check(station string, fail error, finished chan<- string) func(context.Context) error {
	return func(ctx context.Context) error {
		defer func() { finished <- station }()
		if fail != nil {
			return fail
		}
		<-ctx.Done()
		return ctx.Err()
	}
}

// runAllWithin calls runAll and stops the test, saying what to do, if it hasn't returned
// within 1 s.
func runAllWithin(t *testing.T, ctx context.Context, checks []func(context.Context) error, todo string) error {
	t.Helper()
	result := make(chan error, 1)
	go func() {
		result <- runAll(ctx, checks)
	}()
	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		t.Fatalf("runAll(...) is still waiting after 1s: %s", todo)
		return nil
	}
}

func TestRunAllPasses(t *testing.T) {
	ready := func(ctx context.Context) error { return nil }
	err := runAllWithin(t, context.Background(), []func(context.Context) error{ready, ready, ready}, "wait only for the checks")
	if err != nil {
		t.Errorf("runAll(three passing checks) = %v, want nil", err)
	}
}

func TestFirstFailureCancelsTheRest(t *testing.T) {
	errStowed := errors.New("YST: wind 81 km/h, dish stowed")
	finished := make(chan string, 5)
	checks := []func(context.Context) error{
		check("PEL", nil, finished),
		check("SQT", nil, finished),
		check("YST", errStowed, finished),
		check("LAN", nil, finished),
		check("CWR", nil, finished),
	}
	err := runAllWithin(t, context.Background(), checks, "when one check fails, cancel the others")
	if !errors.Is(err, errStowed) {
		t.Errorf("runAll(YST failing, four waiting) = %v, want %v", err, errStowed)
	}
	if n := len(finished); n != 5 {
		t.Errorf("runAll(...) returned while %d of 5 checks were still running, want all of them finished", 5-n)
	}
}

func TestCanceledShiftStopsEveryCheck(t *testing.T) {
	shift, endShift := context.WithCancel(context.Background())
	endShift() // the whole shift is called off
	finished := make(chan string, 2)
	checks := []func(context.Context) error{check("PEL", nil, finished), check("SQT", nil, finished)}
	err := runAllWithin(t, shift, checks, "make the checks' context a child of ctx, so it ends when ctx does")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("runAll(canceled ctx, two waiting checks) = %v, want context.Canceled", err)
	}
	if n := len(finished); n != 2 {
		t.Errorf("runAll(canceled ctx, ...) returned while %d of 2 checks were still running, want both finished", 2-n)
	}
}
