package main

import (
	"errors"
	"testing"
)

func TestRunDecoder(t *testing.T) {
	for _, n := range []int{0, 4182} {
		got, err := runDecoder(func() int { return n })
		if got != n || err != nil {
			t.Errorf("runDecoder(success) = %d, %v, want %d, nil", got, err, n)
		}
	}
	for _, problem := range []any{"broken invariant", errors.New("invalid state")} {
		cleaned := false
		var got int
		var err error
		escaped := catchPanic(func() {
			got, err = runDecoder(func() int { defer func() { cleaned = true }(); panic(problem) })
		})
		if escaped != nil {
			t.Errorf("runDecoder(panic) panicked with %v, want 0, decoder panic: %v", escaped, problem)
			continue
		}
		if got != 0 || err == nil || err.Error() != "decoder panic: "+problemString(problem) {
			t.Errorf("runDecoder(panic) = %d, %v, want 0, decoder panic: %v", got, err, problem)
		}
		if !cleaned {
			t.Error("decoder cleanup = false, want true")
		}
	}
}

// catchPanic runs f and returns the value of a panic that escapes it, or nil.
func catchPanic(f func()) (escaped any) {
	defer func() { escaped = recover() }()
	f()
	return nil
}

// problemString makes the two test panic values comparable with the error text.
func problemString(problem any) string {
	if err, ok := problem.(error); ok {
		return err.Error()
	}
	return problem.(string)
}
