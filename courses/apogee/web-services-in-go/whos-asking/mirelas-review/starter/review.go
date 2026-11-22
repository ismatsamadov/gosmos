package main

import (
	"testing"
	"time"
)

// review checks, over HTTP, that the service at base answers each customer about its own bookings only.
func review(t *testing.T, base string, advance func(time.Duration)) {
	// TODO: check Mirela's points against the service at base, with t.Errorf for every wrong answer.
}
