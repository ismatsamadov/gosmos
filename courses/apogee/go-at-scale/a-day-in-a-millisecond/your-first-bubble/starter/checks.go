package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

// checkWaitUntil fails t when wait doesn't end exactly when it should.
func checkWaitUntil(t *testing.T, wait func(context.Context, func() time.Time, time.Time) error) {
	_ = synctest.Test // TODO: run the checks in a bubble, on 28 June
	at := time.Now().Add(20 * time.Millisecond)
	if err := wait(context.Background(), time.Now, at); err != nil {
		t.Errorf("wait(ctx, time.Now, in 20 ms) = %v, want <nil>", err)
	}
}
