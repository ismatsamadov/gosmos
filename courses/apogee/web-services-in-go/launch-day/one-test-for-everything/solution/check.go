package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// checkService checks the booking service at base over HTTP, the way a customer's program uses it.
func checkService(t *testing.T, base string) {
	if a := request(t, base, "GET", "/health", "", nil); a.Status != http.StatusOK {
		t.Errorf("GET /health = %d, want 200", a.Status)
	}
	furrowlight, brinewatch := token(t, base, furrowlightKey), token(t, base, brinewatchKey)

	var passes []PassView
	a := request(t, base, "GET", "/passes", furrowlight, nil)
	if err := json.Unmarshal(a.Body, &passes); a.Status != http.StatusOK || err != nil || len(passes) == 0 {
		t.Fatalf("GET /passes = %d %s, want 200 and a week of passes", a.Status, a.Body)
	}
	p := passes[0]
	b := Booking{SatID: p.SatID, Station: p.Station, AOS: p.AOS, LOS: p.LOS}
	first := request(t, base, "POST", "/bookings", furrowlight, b, "Idempotency-Key", "rehearsal-1")
	again := request(t, base, "POST", "/bookings", furrowlight, b, "Idempotency-Key", "rehearsal-1")
	if first.Status != http.StatusCreated || again.Status != first.Status || !bytes.Equal(again.Body, first.Body) {
		t.Fatalf("POST /bookings twice with one key = %d, then %d %s, want 201 and the same booking twice", first.Status, again.Status, again.Body)
	}

	json.Unmarshal(first.Body, &b)
	if a := request(t, base, "GET", "/bookings/"+b.ID, brinewatch, nil); a.Status != http.StatusNotFound || a.errorCode() != "not_found" {
		t.Errorf("Brinewatch's GET /bookings/%s = %d %q, want 404 not_found", b.ID, a.Status, a.errorCode())
	}
	for range 11 {
		if a = request(t, base, "GET", "/bookings", brinewatch, nil); a.Status == http.StatusTooManyRequests {
			break
		}
	}
	if a.Status != http.StatusTooManyRequests || a.errorCode() != "rate_limited" {
		t.Errorf("Brinewatch's GET /bookings 11 more times at once = %d %q, want a 429 rate_limited among them", a.Status, a.errorCode())
	}
}
