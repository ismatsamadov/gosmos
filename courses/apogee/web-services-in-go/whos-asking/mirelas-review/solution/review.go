package main

import (
	"net/http"
	"testing"
	"time"
)

// review checks, over HTTP, that the service at base answers each customer about its own bookings only.
func review(t *testing.T, base string, advance func(time.Duration)) {
	for _, route := range privateRoutes {
		if status, _ := call(t, base, route, ""); status != http.StatusUnauthorized {
			t.Errorf("%s with no credentials = %d, want 401", route, status)
		}
	}
	furrowlight := token(t, base, furrowlightKey)
	for _, b := range bookingsOf(t, base, furrowlight) {
		if b.SatID != "furrow-a" && b.SatID != "furrow-b" {
			t.Errorf("GET /bookings as Furrowlight lists %s, a booking of %s", b.ID, b.SatID)
		}
	}
	if status, _ := call(t, base, "GET /bookings", "Bearer "+furrowlightKey); status != http.StatusUnauthorized {
		t.Errorf("GET /bookings with Furrowlight's API key as a token = %d, want 401", status)
	}
	brinewatch := token(t, base, brinewatchKey)
	if status, _ := call(t, base, "GET /bookings/bk-4510", "Bearer "+brinewatch); status != http.StatusNotFound {
		t.Errorf("GET /bookings/bk-4510 as Brinewatch = %d, want 404", status)
	}
	advance(time.Hour)
	status, msg := call(t, base, "GET /bookings", "Bearer "+furrowlight)
	if status != http.StatusUnauthorized || msg != "token expired" {
		t.Errorf("GET /bookings with Furrowlight's token an hour later = %d %q, want 401 %q", status, msg, "token expired")
	}
}
