package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"
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

// Set moves the clock to t.
func (c *clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// Add moves the clock forward by d.
func (c *clock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// received is one request a receiver got, and what it answered.
type received struct {
	Body      string
	Timestamp string
	Status    int
}

// receiver is a customer's webhook receiver, as Yusra's sample writes it. It checks every signature,
// answers 503 to the first refuse requests and 204 after them, and passes on every event it accepts.
type receiver struct {
	*httptest.Server
	mu       sync.Mutex
	refuse   int
	got      []received
	accepted chan Event
}

// newReceiver starts a receiver that checks signatures with secret on the clock now.
func newReceiver(t *testing.T, secret string, now func() time.Time, refuse int) *receiver {
	rc := &receiver{refuse: refuse, accepted: make(chan Event, 10)}
	rc.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))
		ev, err := verifyEvent(r, []byte(secret), now())
		if err != nil {
			t.Errorf("receiver: verifyEvent = %v, want an event signed with %q", err, secret)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		rc.mu.Lock()
		status := http.StatusNoContent
		if rc.refuse > 0 {
			rc.refuse--
			status = http.StatusServiceUnavailable
		}
		rc.got = append(rc.got, received{Body: string(body), Timestamp: r.Header.Get("Apogee-Timestamp"), Status: status})
		rc.mu.Unlock()
		w.WriteHeader(status)
		if status == http.StatusNoContent {
			rc.accepted <- ev
		}
	}))
	t.Cleanup(rc.Close)
	return rc
}

// requests returns every request the receiver got so far.
func (rc *receiver) requests() []received {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return append([]received(nil), rc.got...)
}

// within waits up to two seconds for a value from ch, and fails the test if none comes.
func within[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(2 * time.Second):
		t.Fatalf("waited 2s for %s", what)
		var zero T
		return zero
	}
}

// listen returns a listener on 127.0.0.1 on a port the system picks.
func listen(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln
}

// send makes one request with auth as its Authorization header and body as JSON, and returns the status and body.
func send(method, url, auth string, body any, header ...string) (int, string, error) {
	var r io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return 0, "", err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return resp.StatusCode, string(bytes.TrimSpace(data)), err
}

// bearer gets a token for apiKey from the service at base and returns "Bearer <token>".
func bearer(t *testing.T, base, apiKey string) string {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, base+"/tokens", nil)
	req.Header.Set("Authorization", "Key "+apiKey)
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		t.Fatalf("POST /tokens: %v", err)
	}
	defer resp.Body.Close()
	var tok TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); resp.StatusCode != http.StatusOK || err != nil {
		t.Fatalf("POST /tokens = %s, want 200 and a token", resp.Status)
	}
	return "Bearer " + tok.Token
}

// seed writes launch morning's data folder: four bookings, and an outbox after the beta's 116 events.
func seed(t *testing.T, dir string) {
	t.Helper()
	bookings := []Booking{
		{ID: "bk-4907", SatID: "cirrine-1", Station: "YST", AOS: at("2028-02-01T13:30:16Z"), LOS: at("2028-02-01T13:40:20Z")},
		{ID: "bk-4909", SatID: "brine-3", Station: "LAN", AOS: at("2028-02-01T12:57:33Z"), LOS: at("2028-02-01T13:03:45Z")},
		{ID: "bk-4911", SatID: "furrow-b", Station: "SQT", AOS: at("2028-02-01T14:49:30Z"), LOS: at("2028-02-01T14:58:58Z")},
		{ID: "bk-4912", SatID: "furrow-a", Station: "SQT", AOS: at("2028-02-01T14:02:10Z"), LOS: at("2028-02-01T14:11:52Z"), Note: "launch day"},
	}
	data, err := json.Marshal(bookings)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bookings.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "outbox.json"), []byte(`{"next_event":117,"pending":[]}`), 0600); err != nil {
		t.Fatal(err)
	}
}

// savedOutbox reads outbox.json from dir.
func savedOutbox(t *testing.T, dir string) outboxFile {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "outbox.json"))
	if err != nil {
		t.Fatal(err)
	}
	var o outboxFile
	if err := json.Unmarshal(data, &o); err != nil {
		t.Fatalf("outbox.json = %s: %v", data, err)
	}
	return o
}

// hold holds every request that carries an X-Hold header until release is closed, after telling entered.
func hold(h http.Handler, entered chan<- struct{}, release <-chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Hold") != "" {
			entered <- struct{}{}
			<-release
		}
		h.ServeHTTP(w, r)
	})
}

func TestOpenForBusiness(t *testing.T) {
	dir := t.TempDir()
	seed(t, dir)
	clk := &clock{now: at("2028-02-01T14:04:57Z")}
	cirrine := newReceiver(t, "cirrine webhook secret", clk.Now, 2) // down for Cirrine's own deploy
	furrowlight := newReceiver(t, "furrowlight webhook secret", clk.Now, 0)
	customers := testCustomers()
	for i, c := range customers {
		switch c.ID {
		case "cirrine":
			customers[i].WebhookURL = cirrine.URL
		case "furrowlight":
			customers[i].WebhookURL = furrowlight.URL
		}
	}
	config := Config{
		Dir: dir, Stations: testStations(), Catalog: testCatalog(), Guide: fstest.MapFS{},
		Customers: customers, TokenKey: []byte(testTokenKey), Now: clk.Now,
	}
	limits := Timeouts{Drain: 200 * time.Millisecond, Shutdown: time.Second}

	ok := t.Run("the first process, until 14:05", func(t *testing.T) {
		config.Rand = rand.New(rand.NewPCG(1, 2))
		app, err := newApp(config)
		if err != nil {
			t.Fatalf("newApp = %v, want nil", err)
		}
		secondWait := make(chan struct{}, 1)
		waits := 0
		app.Notifier.Sender.Wait = func(ctx context.Context, d time.Duration) error {
			waits++
			if waits == 1 {
				clk.Add(d)
				return nil
			}
			secondWait <- struct{}{}
			<-ctx.Done() // Cirrine's event is still waiting for its third attempt at 14:05
			return ctx.Err()
		}
		entered, release := make(chan struct{}, 1), make(chan struct{})
		app.Handler = hold(app.Handler, entered, release)
		ln := listen(t)
		base := "http://" + ln.Addr().String()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan error, 1)
		go func() { done <- run(ctx, ln, app, limits) }()

		tomorrow := Booking{SatID: "brine-3", Station: "LAN", AOS: at("2028-02-02T12:27:50Z"), LOS: at("2028-02-02T12:34:37Z")}
		if status, body, err := send(http.MethodPost, base+"/bookings", bearer(t, base, brinewatchKey), tomorrow); status != http.StatusCreated || !bytes.Contains([]byte(body), []byte(`"bk-4913"`)) {
			t.Fatalf("Brinewatch's POST /bookings = %d %s, %v, want 201 and bk-4913", status, body, err)
		}
		// Cirrine-1's file from its 13:30 pass, down at Ystmark by 13:40, crosses the station's slow
		// satellite backhaul and lands at 14:04:57.
		landing := Landing{BookingID: "bk-4907", File: "cirrine-1-20280201-1330.tif", Bytes: 745000000, At: at("2028-02-01T14:04:57Z")}
		if err := app.Notifier.Landed(context.Background(), landing); err != nil {
			t.Fatalf("Landed(bk-4907) = %v, want nil", err)
		}
		within(t, secondWait, "Cirrine's event to be refused twice")

		clk.Set(at("2028-02-01T14:05:00Z"))
		cirrineAuth := bearer(t, base, cirrineKey)
		inFlight := make(chan int, 1)
		go func() {
			status, _, _ := send(http.MethodGet, base+"/bookings", cirrineAuth, nil, "X-Hold", "1")
			inFlight <- status
		}()
		within(t, entered, "Cirrine's GET /bookings to be in flight")
		cancel() // 14:05: the deploy tool stops the old process
		for start := time.Now(); ; {
			status, _, _ := send(http.MethodGet, base+"/health", "", nil)
			if status == http.StatusServiceUnavailable {
				break
			}
			if time.Since(start) > time.Second {
				t.Fatalf("GET /health after the stop = %d for a second, want 503 while the service drains", status)
			}
		}
		close(release)
		if status := within(t, inFlight, "the request in flight to finish"); status != http.StatusOK {
			t.Errorf("Cirrine's GET /bookings, in flight at the stop = %d, want 200", status)
		}
		if err := within(t, done, "run to return"); err != nil {
			t.Errorf("run = %v, want nil", err)
		}
		saved := savedOutbox(t, dir)
		if saved.NextEvent != 118 || len(saved.Pending) != 1 || saved.Pending[0].Event.ID != "evt-000117" || saved.Pending[0].Attempts != 2 {
			t.Errorf("outbox.json after the stop = %+v, want next_event 118 and evt-000117 pending with its 2 attempts", saved)
		}
	})
	if !ok {
		return
	}

	t.Run("the second process, from 14:05:20", func(t *testing.T) {
		clk.Set(at("2028-02-01T14:05:20Z"))
		config.Rand = rand.New(rand.NewPCG(1, 2))
		app, err := newApp(config)
		if err != nil {
			t.Fatalf("newApp(the same folder) = %v, want nil", err)
		}
		app.Notifier.Sender.Wait = func(ctx context.Context, d time.Duration) error {
			clk.Add(d)
			return ctx.Err()
		}
		ln := listen(t)
		base := "http://" + ln.Addr().String()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan error, 1)
		go func() { done <- run(ctx, ln, app, limits) }()

		if ev := within(t, cirrine.accepted, "Cirrine's event after the restart"); ev.ID != "evt-000117" {
			t.Errorf("Cirrine's event after the restart = %q, want evt-000117, the same event", ev.ID)
		}
		if got := len(cirrine.requests()); got != 3 {
			t.Errorf("Cirrine's receiver got %d requests, want 3: two before the restart, one after", got)
		}

		var list []Booking
		status, body, err := send(http.MethodGet, base+"/bookings", bearer(t, base, brinewatchKey), nil)
		if err == nil {
			err = json.Unmarshal([]byte(body), &list)
		}
		if status != http.StatusOK || err != nil || len(list) != 2 || list[1].ID != "bk-4913" {
			t.Errorf("Brinewatch's GET /bookings after the restart = %d %s, want bk-4909 and bk-4913", status, body)
		}
		tomorrow := Booking{SatID: "cirrine-1", Station: "YST", AOS: at("2028-02-02T13:38:14Z"), LOS: at("2028-02-02T13:48:18Z")}
		if status, body, err := send(http.MethodPost, base+"/bookings", bearer(t, base, cirrineKey), tomorrow); status != http.StatusCreated || !bytes.Contains([]byte(body), []byte(`"bk-4914"`)) {
			t.Errorf("Cirrine's POST /bookings after the restart = %d %s, %v, want 201 and bk-4914", status, body, err)
		}

		clk.Set(at("2028-02-01T14:13:03Z"))
		landing := Landing{BookingID: "bk-4912", File: "furrow-a-20280201-1402.tif", Bytes: 10912500000, At: at("2028-02-01T14:13:02Z")}
		if err := app.Notifier.Landed(context.Background(), landing); err != nil {
			t.Fatalf("Landed(bk-4912) = %v, want nil", err)
		}
		within(t, furrowlight.accepted, "Furrow-A's event at Furrowlight's receiver")
		want := `{"id":"evt-000118","type":"data.landed","booking_id":"bk-4912","sat_id":"furrow-a","station":"SQT","file":"furrow-a-20280201-1402.tif","bytes":10912500000,"landed_at":"2028-02-01T14:13:02Z"}`
		if got := furrowlight.requests(); len(got) != 1 || got[0].Body != want || got[0].Timestamp != "1833027183" {
			t.Errorf("Furrowlight's receiver got %+v, want one request at 1833027183, 14:13:03, carrying %s", got, want)
		}
		cancel()
		if err := within(t, done, "run to return"); err != nil {
			t.Errorf("run = %v, want nil", err)
		}
		if saved := savedOutbox(t, dir); saved.NextEvent != 119 || len(saved.Pending) != 0 {
			t.Errorf("outbox.json at the end = %+v, want next_event 119 and nothing pending", saved)
		}
	})
}
