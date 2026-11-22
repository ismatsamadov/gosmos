package main

import (
	"context"
	"errors"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strconv"
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

// cirrineEvent is the event for Cirrine-1's file from Thursday morning's pass over Ystmark.
var cirrineEvent = Event{
	ID:        "evt-000031",
	Type:      "data.landed",
	BookingID: "bk-4752",
	SatID:     "cirrine-1",
	Station:   "YST",
	File:      "cirrine-1-20280120-0843.tif",
	Bytes:     760000000,
	LandedAt:  time.Date(2028, 1, 20, 8, 54, 40, 0, time.UTC),
}

// start is when the first attempt goes out: Thursday 20 January, 08:54:41 UTC.
var start = time.Date(2028, 1, 20, 8, 54, 41, 0, time.UTC)

// attempt is one request a receiver got: the event's ID and the attempt's timestamp.
type attempt struct {
	eventID   string
	timestamp int64
}

// receiver is Cirrine's webhook receiver. It checks every signature as Yusra's sample does, then
// answers with the next status of its script, or 204 once the script is used up.
type receiver struct {
	*httptest.Server
	mu       sync.Mutex
	statuses []int
	got      []attempt
}

// newReceiver starts a receiver that answers with statuses, then 204, checking signatures on clk.
func newReceiver(t *testing.T, clk *clock, statuses ...int) *receiver {
	rc := &receiver{statuses: statuses}
	rc.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ev, err := verifyEvent(r, []byte("cirrine webhook secret"), clk.Now())
		if err != nil {
			t.Errorf("Cirrine's receiver: verifyEvent = %v, want a signed event", err)
		}
		rc.mu.Lock()
		status := http.StatusNoContent
		if len(rc.statuses) > 0 {
			status, rc.statuses = rc.statuses[0], rc.statuses[1:]
		}
		ts, _ := strconv.ParseInt(r.Header.Get("Apogee-Timestamp"), 10, 64)
		rc.got = append(rc.got, attempt{ev.ID, ts})
		rc.mu.Unlock()
		w.WriteHeader(status)
	}))
	t.Cleanup(rc.Close)
	return rc
}

// attempts returns every request the receiver got so far.
func (rc *receiver) attempts() []attempt {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return append([]attempt(nil), rc.got...)
}

// newSender returns a Sender for Cirrine's receiver at url whose waits move clk and are kept in waits.
func newSender(url string, clk *clock, waits *[]time.Duration) *Sender {
	return &Sender{
		Client:    &http.Client{Timeout: 2 * time.Second},
		Customers: map[string]Customer{"cirrine": {ID: "cirrine", WebhookURL: url, WebhookKey: []byte("cirrine webhook secret")}},
		Now:       clk.Now,
		Rand:      rand.New(rand.NewPCG(1, 2)),
		Wait: func(ctx context.Context, d time.Duration) error {
			*waits = append(*waits, d)
			clk.Add(d)
			return ctx.Err()
		},
	}
}

func TestDeliverRetriesUntilItArrives(t *testing.T) {
	clk := &clock{now: start}
	rc := newReceiver(t, clk, http.StatusServiceUnavailable, http.StatusServiceUnavailable)
	var waits []time.Duration
	d := &Delivery{Customer: "cirrine", Event: cirrineEvent}
	if err := newSender(rc.URL, clk, &waits).Deliver(context.Background(), d); err != nil {
		t.Fatalf("Deliver(503, 503, 204) = %v, want nil", err)
	}
	if d.Attempts != 3 {
		t.Errorf("Deliver(503, 503, 204) left Attempts = %d, want 3", d.Attempts)
	}
	if len(waits) != 2 || waits[0] < time.Second || waits[0] >= 2*time.Second || waits[1] < 2*time.Second || waits[1] >= 4*time.Second {
		t.Fatalf("Deliver(503, 503, 204) waited %v, want one wait from 1s up to 2s, then one from 2s up to 4s", waits)
	}
	got := rc.attempts()
	want := []attempt{
		{"evt-000031", start.Unix()},
		{"evt-000031", start.Add(waits[0]).Unix()},
		{"evt-000031", start.Add(waits[0] + waits[1]).Unix()},
	}
	if len(got) != len(want) {
		t.Fatalf("Cirrine's receiver got %d requests, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("attempt %d reached the receiver as %+v, want %+v: the same event, timestamped when it was sent", i+1, got[i], want[i])
		}
	}
}

func TestDeliverGivesUp(t *testing.T) {
	tests := []struct {
		name     string
		statuses []int
		attempts int
	}{
		{"a receiver down all day", []int{503, 503, 503, 503, 503, 503, 503, 503, 503}, 8},
		{"an address nobody serves", []int{http.StatusNotFound}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := &clock{now: start}
			rc := newReceiver(t, clk, tt.statuses...)
			var waits []time.Duration
			d := &Delivery{Customer: "cirrine", Event: cirrineEvent}
			err := newSender(rc.URL, clk, &waits).Deliver(context.Background(), d)
			if !errors.Is(err, ErrUndelivered) {
				t.Errorf("Deliver(%s) = %v, want ErrUndelivered", tt.name, err)
			}
			if d.Attempts != tt.attempts || len(rc.attempts()) != tt.attempts {
				t.Errorf("Deliver(%s) left Attempts = %d and the receiver got %d requests, want %d", tt.name, d.Attempts, len(rc.attempts()), tt.attempts)
			}
			if len(waits) != tt.attempts-1 {
				t.Errorf("Deliver(%s) waited %d times, want %d: no wait after the last attempt", tt.name, len(waits), tt.attempts-1)
			}
		})
	}
}

func TestDeliverWithNoReceiver(t *testing.T) {
	clk := &clock{now: start}
	gone := httptest.NewServer(http.NotFoundHandler())
	gone.Close()
	var waits []time.Duration
	d := &Delivery{Customer: "cirrine", Event: cirrineEvent}
	err := newSender(gone.URL, clk, &waits).Deliver(context.Background(), d)
	if !errors.Is(err, ErrUndelivered) || d.Attempts != 8 {
		t.Errorf("Deliver(to a refused connection) = %v with Attempts = %d, want ErrUndelivered with 8: no answer at all is worth retrying", err, d.Attempts)
	}
}

func TestDeliverStopsWithTheService(t *testing.T) {
	clk := &clock{now: start}
	rc := newReceiver(t, clk, 503, 503, 503)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var waits []time.Duration
	s := newSender(rc.URL, clk, &waits)
	s.Wait = func(ctx context.Context, d time.Duration) error {
		waits = append(waits, d)
		clk.Add(d)
		if len(waits) == 2 {
			cancel() // the service stops during the second wait
		}
		return ctx.Err()
	}
	d := &Delivery{Customer: "cirrine", Event: cirrineEvent}
	err := s.Deliver(ctx, d)
	if !errors.Is(err, context.Canceled) || errors.Is(err, ErrUndelivered) {
		t.Errorf("Deliver(stopped during a wait) = %v, want context.Canceled and not ErrUndelivered: stopping isn't failing", err)
	}
	if d.Attempts != 2 {
		t.Errorf("Deliver(stopped during a wait) left Attempts = %d, want 2, the attempts made", d.Attempts)
	}
	if !d.NextAt.Equal(clk.Now()) {
		t.Errorf("Deliver(stopped during a wait) left NextAt = %v, want %v, when the wait would have ended", d.NextAt, clk.Now())
	}
}

func TestDeliverStopsDuringAnAttempt(t *testing.T) {
	arrived, release := make(chan struct{}, 1), make(chan struct{})
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.ReadAll(r.Body)
		arrived <- struct{}{}
		select { // Cirrine's receiver is still thinking when the service stops
		case <-r.Context().Done():
		case <-release:
		}
	}))
	t.Cleanup(slow.Close)
	t.Cleanup(func() { close(release) })
	clk := &clock{now: start}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var waits []time.Duration
	d := &Delivery{Customer: "cirrine", Event: cirrineEvent}
	done := make(chan error, 1)
	go func() { done <- newSender(slow.URL, clk, &waits).Deliver(ctx, d) }()
	select {
	case <-arrived:
	case <-time.After(time.Second):
		t.Fatal("the first attempt never reached Cirrine's receiver")
	}
	cancel() // the service stops while the attempt is still waiting for its answer
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) || errors.Is(err, ErrUndelivered) {
			t.Errorf("Deliver(stopped during an attempt) = %v, want context.Canceled and not ErrUndelivered: stopping isn't failing", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Deliver hadn't returned a second after its context ended")
	}
}
