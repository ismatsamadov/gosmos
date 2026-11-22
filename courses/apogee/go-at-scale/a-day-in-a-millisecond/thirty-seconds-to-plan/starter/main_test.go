package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"testing/synctest"
	"time"
)

// version writes a schedule as its version, or <nil>.
func version(s *Schedule) string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("version %d", s.Version)
}

// sentinel names one of context's two errors as Go code writes it.
func sentinel(err error) string {
	if err == context.Canceled {
		return "context.Canceled"
	}
	return "context.DeadlineExceeded"
}

// errText writes an error's message in quotes, or <nil>.
func errText(err error) string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%q", err.Error())
}

func TestPlanWithin(t *testing.T) {
	run41 := &Schedule{Version: 41}
	tests := []struct {
		name     string
		limit    time.Duration // planWithin's limit
		run      time.Duration // how long the planning run takes, watching its context
		fail     error         // what the run returns once done, instead of run41
		cancelAt time.Duration // when the caller cancels; 0 for never
		want     *Schedule
		wantErr  string // the error's message; "" for nil
		wantIs   error  // what errors.Is must find in the error
		wantAt   time.Duration
	}{
		{"a 29 s run", 30 * time.Second, 29 * time.Second, nil, 0, run41, "", nil, 29 * time.Second},
		{"a 31 s run", 30 * time.Second, 31 * time.Second, nil, 0, nil, "plan: took longer than 30s: context deadline exceeded", context.DeadlineExceeded, 30 * time.Second},
		{"a 31 s run with a 20 s limit", 20 * time.Second, 31 * time.Second, nil, 0, nil, "plan: took longer than 20s: context deadline exceeded", context.DeadlineExceeded, 20 * time.Second},
		{"a run its caller cancels at 10 s", 30 * time.Second, 31 * time.Second, nil, 10 * time.Second, nil, "context canceled", context.Canceled, 10 * time.Second},
		{"a run that fails at 5 s", 30 * time.Second, 5 * time.Second, errors.New("predict cirrine-7: no elements"), 0, nil, "predict cirrine-7: no elements", nil, 5 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				start := time.Now()
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				if tt.cancelAt > 0 {
					time.AfterFunc(tt.cancelAt, cancel)
				}
				var runCtx context.Context // the context the run was given
				got, err := planWithin(ctx, tt.limit, func(ctx context.Context) (*Schedule, error) {
					runCtx = ctx
					select {
					case <-time.After(tt.run):
						if tt.fail != nil {
							return nil, tt.fail
						}
						return run41, nil
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				})
				at := time.Since(start)
				want := "<nil>"
				if tt.wantErr != "" {
					want = fmt.Sprintf("%q", tt.wantErr)
				}
				if got != tt.want || errText(err) != want || at != tt.wantAt {
					t.Errorf("planWithin(ctx, %v, %s) = %s, %s at %v, want %s, %s at %v",
						tt.limit, tt.name, version(got), errText(err), at, version(tt.want), want, tt.wantAt)
				}
				if err != nil && tt.wantIs != nil && !errors.Is(err, tt.wantIs) {
					t.Errorf("planWithin(ctx, %v, %s): errors.Is(err, %s) = false, want true", tt.limit, tt.name, sentinel(tt.wantIs))
				}
				if runCtx != nil && runCtx.Err() == nil {
					t.Errorf("planWithin(ctx, %v, %s) returned with its run's context still live, want it canceled", tt.limit, tt.name)
				}
			})
		})
	}
}
