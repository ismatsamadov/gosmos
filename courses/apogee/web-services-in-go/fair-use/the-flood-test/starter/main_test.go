package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"sync"
	"testing"
	"testing/fstest"
	"time"
)

// tokenKey is the test token key: 32 bytes, like the service's real key, which never appears in a lesson.
const tokenKey = "bookingd test token key 32 bytes"

// clock is the service's clock in these tests: it stands still unless a test moves it.
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

// at parses an RFC 3339 time from the fixtures.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// bookings returns a store holding this week's bookings of Brinewatch's and Cirrine's.
func bookings(t *testing.T) *MemStore {
	t.Helper()
	store := NewMemStore()
	err := store.load([]Booking{
		{ID: "bk-4520", SatID: "brine-3", Station: "LAN", AOS: at("2028-01-15T06:17:08Z"), LOS: at("2028-01-15T06:24:30Z")},
		{ID: "bk-4521", SatID: "cirrine-1", Station: "YST", AOS: at("2028-01-15T08:03:29Z"), LOS: at("2028-01-15T08:13:43Z")},
	})
	if err != nil {
		t.Fatal(err)
	}
	return store
}

// newHandler returns the whole service, middleware and all, as it runs on Friday 14 January at 10:00 UTC.
func newHandler(store Store, timeouts Timeouts) (http.Handler, *clock) {
	c := &clock{now: time.Date(2028, 1, 14, 10, 0, 0, 0, time.UTC)}
	s := &Service{
		Stations: []Station{{Code: "YST", Name: "Ystmark", LatDeg: 77.8, LonDeg: 2.0, MaskDeg: 5, Bands: []Band{"S", "X"}}},
		Catalog: []PassView{
			{ID: "brine-3/LAN/09824", Pass: Pass{SatID: "brine-3", Station: "LAN", AOS: at("2028-01-15T06:17:08Z"), LOS: at("2028-01-15T06:24:30Z"), MaxElevDeg: 29.8}},
			{ID: "cirrine-1/YST/11961", Pass: Pass{SatID: "cirrine-1", Station: "YST", AOS: at("2028-01-15T08:03:29Z"), LOS: at("2028-01-15T08:13:43Z"), MaxElevDeg: 56.7}},
		},
		Guide: fstest.MapFS{"index.html": {Data: []byte("<h1>Booking passes with Apogee</h1>\n")}},
		Store: store,
		Customers: map[string]Customer{
			"brinewatch": {ID: "brinewatch", Name: "Brinewatch", SatIDs: []string{"brine-3"}, KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
			"cirrine":    {ID: "cirrine", Name: "Cirrine Atmospherics", SatIDs: []string{"cirrine-1"}, KeyHash: "d9ee2e5de1796bbeba223c0737933b3e60d53687d479f8a04c98c99df87b510d"},
		},
		TokenKey: []byte(tokenKey),
		Limiter:  newLimiter(c.Now),
		Timeouts: timeouts,
		Now:      c.Now,
	}
	return s.handler(log.New(io.Discard, "", 0)), c
}

// bearer returns the Authorization header of a token for customer that is good until 11:00.
func bearer(customer string) string {
	claims := customer + "." + strconv.FormatInt(time.Date(2028, 1, 14, 11, 0, 0, 0, time.UTC).Unix(), 10)
	mac := hmac.New(sha256.New, []byte(tokenKey))
	mac.Write([]byte(claims))
	return "Bearer " + claims + "." + hex.EncodeToString(mac.Sum(nil))
}

// get sends GET path to h with auth as its Authorization header ("" sends none).
func get(h http.Handler, path, auth string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, path, nil)
	if auth != "" {
		r.Header.Set("Authorization", auth)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

// count tallies how many answers had each status, as "200: 10, 429: 990".
func count(codes []int) string {
	tally := map[int]int{}
	for _, code := range codes {
		tally[code]++
	}
	keys := []int{}
	for code := range tally {
		keys = append(keys, code)
	}
	slices.Sort(keys)
	out := ""
	for _, code := range keys {
		out += fmt.Sprintf(", %d: %d", code, tally[code])
	}
	if out == "" {
		return "none"
	}
	return out[2:]
}

func TestFlood(t *testing.T) {
	h, _ := newHandler(bookings(t), Timeouts{Handler: 10 * time.Second})
	var brinewatch, cirrine []int
	badRetry := 0
	for i := 1; i <= 1000; i++ {
		rec := get(h, "/bookings", bearer("brinewatch"))
		brinewatch = append(brinewatch, rec.Code)
		if rec.Code == http.StatusTooManyRequests && rec.Header().Get("Retry-After") != "1" {
			badRetry++
		}
		if i%200 == 0 {
			cirrine = append(cirrine, get(h, "/bookings", bearer("cirrine")).Code)
		}
	}
	if got, want := count(brinewatch), "200: 10, 429: 990"; got != want {
		t.Errorf("Brinewatch's 1,000 GET /bookings at 10:00:00 = %s, want %s", got, want)
	}
	if badRetry > 0 {
		t.Errorf("Brinewatch's refused requests without Retry-After: 1 = %d, want 0", badRetry)
	}
	if got, want := count(cirrine), "200: 5"; got != want {
		t.Errorf("Cirrine's 5 GET /bookings during the flood = %s, want %s", got, want)
	}
}

func TestNoTokenSpendsNothing(t *testing.T) {
	h, _ := newHandler(bookings(t), Timeouts{Handler: 10 * time.Second})
	var none []int
	for range 20 {
		rec := get(h, "/bookings", "")
		none = append(none, rec.Code)
		if got := rec.Header().Get("X-RateLimit-Remaining"); got != "" {
			t.Fatalf("GET /bookings with no token: X-RateLimit-Remaining = %q, want none: limit must not see it", got)
		}
	}
	for range 20 {
		none = append(none, get(h, "/bookings", "Bearer brinewatch.1831460400.forged").Code)
	}
	if got, want := count(none), "401: 40"; got != want {
		t.Errorf("40 GET /bookings with no token or a forged one = %s, want %s", got, want)
	}
	var brinewatch []int
	for range 11 {
		brinewatch = append(brinewatch, get(h, "/bookings", bearer("brinewatch")).Code)
	}
	if got, want := count(brinewatch), "200: 10, 429: 1"; got != want {
		t.Errorf("Brinewatch's 11 GET /bookings after 40 refused ones = %s, want %s", got, want)
	}
}

func TestFloodRefills(t *testing.T) {
	h, c := newHandler(bookings(t), Timeouts{Handler: 10 * time.Second})
	for range 1000 {
		get(h, "/bookings", bearer("brinewatch"))
	}
	c.Add(5 * time.Second)
	var later []int
	for range 10 {
		later = append(later, get(h, "/bookings", bearer("brinewatch")).Code)
	}
	if got, want := count(later), "200: 5, 429: 5"; got != want {
		t.Errorf("Brinewatch's 10 GET /bookings five seconds after its flood = %s, want %s", got, want)
	}
}

// stalledStore is the store stuck behind the flood: Get waits until its context ends, or half a second.
type stalledStore struct {
	Store
}

// Get waits for ctx to end and returns why it ended.
func (stalledStore) Get(ctx context.Context, id string) (Booking, error) {
	select {
	case <-ctx.Done():
		return Booking{}, ctx.Err()
	case <-time.After(500 * time.Millisecond):
		return Booking{}, errors.New("the store waited half a second: the request's context had no deadline")
	}
}

func TestDeadlineAroundTheMux(t *testing.T) {
	h, _ := newHandler(stalledStore{bookings(t)}, Timeouts{Handler: 50 * time.Millisecond})
	rec := get(h, "/bookings/bk-4521", bearer("cirrine"))
	var body ErrorBody
	got := fmt.Sprintf("%d %q", rec.Code, rec.Body)
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err == nil {
		got = fmt.Sprintf("%d %s %q", rec.Code, body.Code, body.Message)
	}
	if want := `503 unavailable "request took too long"`; got != want {
		t.Errorf("Cirrine's GET /bookings/bk-4521 with a stalled store, at a 50ms deadline = %s, want %s", got, want)
	}
	if id := rec.Header().Get("X-Request-Id"); id == "" || body.RequestID != id {
		t.Errorf("the 503's request_id = %q, X-Request-Id = %q, want the same ID in both", body.RequestID, id)
	}
}
