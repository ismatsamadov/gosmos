package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDeadlineIsLOS(t *testing.T) {
	los := time.Now().UTC().Add(9 * time.Minute) // the pass ends in nine minutes
	ctx, cancel := passContext(context.Background(), los)
	defer cancel()
	got, ok := ctx.Deadline()
	if !ok || !got.Equal(los) {
		t.Errorf("passContext(ctx, los).Deadline() = %v, %v; want %v, true", got, ok, los)
	}
	if err := ctx.Err(); err != nil {
		t.Errorf("passContext(ctx, los nine minutes away).Err() = %v, want nil", err)
	}
}

func TestPassAlreadyOver(t *testing.T) {
	los := time.Now().UTC().Add(-time.Minute) // LOS was a minute ago
	ctx, cancel := passContext(context.Background(), los)
	defer cancel()
	if err := ctx.Err(); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("passContext(ctx, los a minute ago).Err() = %v, want context.DeadlineExceeded", err)
	}
}

func TestCanceledParentEndsThePass(t *testing.T) {
	shift, endShift := context.WithCancel(context.Background())
	endShift() // the whole shift is called off
	ctx, cancel := passContext(shift, time.Now().UTC().Add(9*time.Minute))
	defer cancel()
	if err := ctx.Err(); !errors.Is(err, context.Canceled) {
		t.Errorf("passContext(canceled parent, los).Err() = %v, want context.Canceled", err)
	}
}
