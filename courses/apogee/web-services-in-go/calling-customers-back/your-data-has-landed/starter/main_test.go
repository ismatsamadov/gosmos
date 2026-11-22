package main

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// clock is a fake clock that moves only when the test moves it.
type clock struct {
	mu  sync.Mutex
	now time.Time
}

// Now returns the clock's time.
func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Add moves the clock forward by d.
func (c *clock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// receiver is a customer's webhook receiver. It checks every signature, answers with the next
// status of its script or 204 once the script is used up, and keeps every event it was sent.
type receiver struct {
	*httptest.Server
	mu       sync.Mutex
	statuses []int
	events   []Event
}

// newReceiver starts a receiver for a customer whose webhook secret is secret, on the clock clk.
func newReceiver(t *testing.T, clk *clock, secret string, statuses ...int) *receiver {
	rc := &receiver{statuses: statuses}
	rc.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ev, err := verifyEvent(r, []byte(secret), clk.Now())
		if err != nil {
			t.Errorf("receiver: verifyEvent = %v, want an event signed with %q", err, secret)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		rc.mu.Lock()
		defer rc.mu.Unlock()
		rc.events = append(rc.events, ev)
		status := http.StatusNoContent
		if len(rc.statuses) > 0 {
			status, rc.statuses = rc.statuses[0], rc.statuses[1:]
		}
		w.WriteHeader(status)
	}))
	t.Cleanup(rc.Close)
	return rc
}

// got returns every event the receiver was sent so far, one per attempt.
func (rc *receiver) got() []Event {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return append([]Event(nil), rc.events...)
}

// friday holds Friday's bookings, as the store keeps them.
var friday = []Booking{
	{ID: "bk-4779", SatID: "brine-3", Station: "LAN", AOS: time.Date(2028, 1, 21, 16, 42, 46, 0, time.UTC), LOS: time.Date(2028, 1, 21, 16, 50, 14, 0, time.UTC)},
	{ID: "bk-4781", SatID: "cirrine-1", Station: "YST", AOS: time.Date(2028, 1, 21, 8, 51, 36, 0, time.UTC), LOS: time.Date(2028, 1, 21, 9, 1, 56, 0, time.UTC)},
	{ID: "bk-4783", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2028, 1, 21, 14, 31, 40, 0, time.UTC), LOS: time.Date(2028, 1, 21, 14, 41, 39, 0, time.UTC)},
}

// The landings: Furrow-A's image at 14:42:38, and Cirrine-1's file from the morning.
var (
	furrowLanding  = Landing{BookingID: "bk-4783", File: "furrow-a-20280121-1431.tif", Bytes: 10893750000, At: time.Date(2028, 1, 21, 14, 42, 38, 0, time.UTC)}
	cirrineLanding = Landing{BookingID: "bk-4781", File: "cirrine-1-20280121-0851.tif", Bytes: 762500000, At: time.Date(2028, 1, 21, 9, 2, 50, 0, time.UTC)}
)

// setup returns a Notifier for Friday's bookings, telling Furrowlight's and Cirrine's receivers, and where its waits go.
func setup(t *testing.T, clk *clock, furrowlight, cirrine *receiver) (*Notifier, *[]time.Duration) {
	t.Helper()
	store := NewMemStore()
	if err := store.load(friday); err != nil {
		t.Fatal(err)
	}
	customers := map[string]Customer{
		"furrowlight": {ID: "furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}, WebhookURL: furrowlight.URL, WebhookKey: []byte("furrowlight webhook secret")},
		"brinewatch":  {ID: "brinewatch", SatIDs: []string{"brine-3"}},
		"cirrine":     {ID: "cirrine", SatIDs: []string{"cirrine-1"}, WebhookURL: cirrine.URL, WebhookKey: []byte("cirrine webhook secret")},
	}
	var waits []time.Duration
	sender := &Sender{
		Client:    &http.Client{Timeout: 2 * time.Second},
		Customers: customers,
		Now:       clk.Now,
		Rand:      rand.New(rand.NewPCG(1, 2)),
		Wait: func(ctx context.Context, d time.Duration) error {
			waits = append(waits, d)
			clk.Add(d)
			return ctx.Err()
		},
	}
	return &Notifier{Store: store, Customers: customers, Sender: sender}, &waits
}

// at1442 is when the landing reaches the notifier: Friday 21 January, 14:42:39 UTC, a minute after LOS.
var at1442 = time.Date(2028, 1, 21, 14, 42, 39, 0, time.UTC)

func TestLandedTellsFurrowlight(t *testing.T) {
	clk := &clock{now: at1442}
	furrowlight := newReceiver(t, clk, "furrowlight webhook secret")
	cirrine := newReceiver(t, clk, "cirrine webhook secret")
	n, _ := setup(t, clk, furrowlight, cirrine)
	if err := n.Landed(context.Background(), furrowLanding); err != nil {
		t.Fatalf("Landed(bk-4783) = %v, want nil", err)
	}
	want := Event{
		ID:        "evt-000001",
		Type:      "data.landed",
		BookingID: "bk-4783",
		SatID:     "furrow-a",
		Station:   "SQT",
		File:      "furrow-a-20280121-1431.tif",
		Bytes:     10893750000,
		LandedAt:  furrowLanding.At,
	}
	got := furrowlight.got()
	if len(got) != 1 || got[0] != want {
		t.Errorf("Furrowlight's receiver got %+v, want one event %+v", got, want)
	}
	if others := cirrine.got(); len(others) > 0 {
		t.Errorf("Cirrine's receiver got %+v, want nothing: Furrow-A is Furrowlight's", others)
	}
}

func TestLandedTellsTheRightCustomer(t *testing.T) {
	clk := &clock{now: time.Date(2028, 1, 21, 9, 2, 51, 0, time.UTC)}
	furrowlight := newReceiver(t, clk, "furrowlight webhook secret")
	cirrine := newReceiver(t, clk, "cirrine webhook secret")
	n, _ := setup(t, clk, furrowlight, cirrine)
	if err := n.Landed(context.Background(), cirrineLanding); err != nil {
		t.Fatalf("Landed(bk-4781) = %v, want nil", err)
	}
	if got := cirrine.got(); len(got) != 1 || got[0].BookingID != "bk-4781" || got[0].SatID != "cirrine-1" {
		t.Errorf("Cirrine's receiver got %+v, want one event for bk-4781", got)
	}
	if others := furrowlight.got(); len(others) > 0 {
		t.Errorf("Furrowlight's receiver got %+v, want nothing: Cirrine-1 is Cirrine's", others)
	}
}

func TestLandedRetries(t *testing.T) {
	clk := &clock{now: at1442}
	furrowlight := newReceiver(t, clk, "furrowlight webhook secret", http.StatusServiceUnavailable)
	n, waits := setup(t, clk, furrowlight, newReceiver(t, clk, "cirrine webhook secret"))
	if err := n.Landed(context.Background(), furrowLanding); err != nil {
		t.Fatalf("Landed(bk-4783, receiver failing once) = %v, want nil", err)
	}
	got := furrowlight.got()
	if len(got) != 2 || got[0].ID != "evt-000001" || got[1].ID != "evt-000001" {
		t.Errorf("Furrowlight's receiver got %+v, want evt-000001 twice: a retry keeps the event's ID", got)
	}
	if len(*waits) != 1 {
		t.Errorf("Landed(bk-4783, receiver failing once) waited %d times, want 1", len(*waits))
	}
}

func TestLandedUnknownBooking(t *testing.T) {
	clk := &clock{now: at1442}
	furrowlight := newReceiver(t, clk, "furrowlight webhook secret")
	n, _ := setup(t, clk, furrowlight, newReceiver(t, clk, "cirrine webhook secret"))
	lost := Landing{BookingID: "bk-5499", File: "furrow-a-20280121-1431.tif", Bytes: 10893750000, At: furrowLanding.At}
	if err := n.Landed(context.Background(), lost); !errors.Is(err, ErrNotFound) {
		t.Errorf("Landed(bk-5499) = %v, want ErrNotFound", err)
	}
	if got := furrowlight.got(); len(got) > 0 {
		t.Errorf("Landed(bk-5499) sent %+v, want nothing sent", got)
	}
}

func TestEventIDsCount(t *testing.T) {
	clk := &clock{now: at1442}
	furrowlight := newReceiver(t, clk, "furrowlight webhook secret")
	cirrine := newReceiver(t, clk, "cirrine webhook secret")
	n, _ := setup(t, clk, furrowlight, cirrine)
	n.Landed(context.Background(), furrowLanding)
	n.Landed(context.Background(), cirrineLanding)
	first, second := furrowlight.got(), cirrine.got()
	if len(first) != 1 || first[0].ID != "evt-000001" || len(second) != 1 || second[0].ID != "evt-000002" {
		t.Errorf("two landings sent %+v and %+v, want evt-000001 and then evt-000002", first, second)
	}
}
