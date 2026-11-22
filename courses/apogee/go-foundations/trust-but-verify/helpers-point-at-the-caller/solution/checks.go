package main

import (
	"math"
	"testing"
)

// assertNear fails the test when got is more than tol away from want.
func assertNear(t testing.TB, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("got %.2f, want %.2f within %.2f", got, want, tol)
	}
}
