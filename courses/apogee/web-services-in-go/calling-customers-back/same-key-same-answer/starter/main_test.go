package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
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

// The bookings the tests send: Friday's passes over Salar Quieto, and one each for Brinewatch and Cirrine.
const (
	furrowA  = `{"sat_id":"furrow-a","station":"SQT","aos":"2028-01-21T14:31:40Z","los":"2028-01-21T14:41:39Z"}`
	furrowB  = `{"sat_id":"furrow-b","station":"SQT","aos":"2028-01-21T15:19:53Z","los":"2028-01-21T15:28:16Z"}`
	brine3   = `{"sat_id":"brine-3","station":"LAN","aos":"2028-01-21T16:42:46Z","los":"2028-01-21T16:50:14Z"}`
	cirrine1 = `{"sat_id":"cirrine-1","station":"YST","aos":"2028-01-21T08:51:36Z","los":"2028-01-21T09:01:56Z"}`
)

// bookPass stands in for createBooking: it keeps the booking it is sent and answers 201 with its Location.
func bookPass(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b Booking
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			writeError(w, r, &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: "request body: " + err.Error()})
			return
		}
		b, err := store.Create(r.Context(), b)
		if err != nil {
			writeError(w, r, err)
			return
		}
		w.Header().Set("Location", "/bookings/"+b.ID)
		writeJSON(w, http.StatusCreated, b)
	}
}

// bookings returns POST /bookings as the service wraps it for a customer: idempotent, then the booking.
func bookings(keys *keyStore, store Store) func(customer string) http.Handler {
	h := idempotent(keys)(bookPass(store))
	return func(customer string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(withCustomer(r.Context(), Customer{ID: customer})))
		})
	}
}

// send makes one POST /bookings through h, with key as its Idempotency-Key unless it is empty.
func send(h http.Handler, key, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, "/bookings", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	if key != "" {
		r.Header.Set("Idempotency-Key", key)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

// count returns how many bookings store holds.
func count(t *testing.T, store Store) int {
	t.Helper()
	list, err := store.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return len(list)
}

// text returns an answer's body without its final newline.
func text(rec *httptest.ResponseRecorder) string {
	return strings.TrimSpace(rec.Body.String())
}

// errorOf decodes the ErrorBody of an answer.
func errorOf(rec *httptest.ResponseRecorder) ErrorBody {
	var body ErrorBody
	json.Unmarshal(rec.Body.Bytes(), &body)
	return body
}

// friday is when the tests start: Friday 21 January, 08:00 UTC.
var friday = time.Date(2028, 1, 21, 8, 0, 0, 0, time.UTC)

func TestRetryGetsTheSameAnswer(t *testing.T) {
	store := NewMemStore()
	furrowlight := bookings(newKeyStore((&clock{now: friday}).Now), store)("furrowlight")
	first := send(furrowlight, "wk03-fri", furrowA)
	if first.Code != http.StatusCreated || first.Header().Get("Location") != "/bookings/bk-4501" {
		t.Fatalf("first POST /bookings = %d at %q, want 201 at /bookings/bk-4501", first.Code, first.Header().Get("Location"))
	}
	if got := first.Header().Get("Idempotent-Replayed"); got != "" {
		t.Errorf("first POST /bookings has Idempotent-Replayed %q, want none: it ran for real", got)
	}
	again := send(furrowlight, "wk03-fri", furrowA)
	if again.Code != first.Code || again.Body.String() != first.Body.String() || again.Header().Get("Location") != "/bookings/bk-4501" {
		t.Errorf("retry with the same key = %d %s at %q, want the first answer again: %d %s at /bookings/bk-4501",
			again.Code, text(again), again.Header().Get("Location"), first.Code, text(first))
	}
	if got := again.Header().Get("Idempotent-Replayed"); got != "true" {
		t.Errorf("retry with the same key has Idempotent-Replayed %q, want %q", got, "true")
	}
	if n := count(t, store); n != 1 {
		t.Errorf("after a retry the store holds %d bookings, want 1", n)
	}
}

func TestKeyWithAnotherBody(t *testing.T) {
	store := NewMemStore()
	furrowlight := bookings(newKeyStore((&clock{now: friday}).Now), store)("furrowlight")
	send(furrowlight, "wk03-fri", furrowA)
	rec := send(furrowlight, "wk03-fri", furrowB)
	want := ErrorBody{Code: "conflict", Message: `idempotency key "wk03-fri" was used for a different request`}
	if got := errorOf(rec); rec.Code != http.StatusConflict || got != want {
		t.Errorf("the same key with another body = %d %+v, want 409 %+v", rec.Code, got, want)
	}
	if n := count(t, store); n != 1 {
		t.Errorf("after a refused request the store holds %d bookings, want 1", n)
	}
}

// slowStore holds the first Create until release is closed, after closing entered.
type slowStore struct {
	Store
	held             atomic.Bool
	entered, release chan struct{}
}

// Create waits for release the first time it is called, then keeps b.
func (s *slowStore) Create(ctx context.Context, b Booking) (Booking, error) {
	if s.held.CompareAndSwap(false, true) {
		close(s.entered)
		<-s.release
	}
	return s.Store.Create(ctx, b)
}

func TestKeyInProgress(t *testing.T) {
	store := &slowStore{Store: NewMemStore(), entered: make(chan struct{}), release: make(chan struct{})}
	furrowlight := bookings(newKeyStore((&clock{now: friday}).Now), store)("furrowlight")
	firstDone := make(chan *httptest.ResponseRecorder, 1)
	go func() { firstDone <- send(furrowlight, "wk03-fri", furrowA) }()
	select {
	case <-store.entered:
	case <-time.After(time.Second):
		t.Fatal("the first request never reached the store")
	}
	rec := send(furrowlight, "wk03-fri", furrowA)
	close(store.release)
	want := ErrorBody{Code: "conflict", Message: `a request with idempotency key "wk03-fri" is in progress`}
	if got := errorOf(rec); rec.Code != http.StatusConflict || got != want {
		t.Errorf("the same request while the first still runs = %d %+v, want 409 %+v", rec.Code, got, want)
	}
	if first := <-firstDone; first.Code != http.StatusCreated {
		t.Errorf("the first request = %d %s, want 201", first.Code, text(first))
	}
}

func TestKeyMustBeWellFormed(t *testing.T) {
	furrowlight := bookings(newKeyStore((&clock{now: friday}).Now), NewMemStore())("furrowlight")
	want := ErrorBody{Code: "bad_request", Message: "idempotency key must be 1 to 64 letters, digits, - or _"}
	for _, key := range []string{"!!", strings.Repeat("k", 65)} {
		if rec := send(furrowlight, key, furrowA); rec.Code != http.StatusBadRequest || errorOf(rec) != want {
			t.Errorf("POST /bookings with key %q = %d %s, want 400 %+v", key, rec.Code, text(rec), want)
		}
	}
	if rec := send(furrowlight, strings.Repeat("k", 64), furrowA); rec.Code != http.StatusCreated {
		t.Errorf("POST /bookings with a 64-character key = %d %s, want 201", rec.Code, text(rec))
	}
}

func TestKeyLastsADay(t *testing.T) {
	clk := &clock{now: friday}
	furrowlight := bookings(newKeyStore(clk.Now), NewMemStore())("furrowlight")
	send(furrowlight, "wk03-fri", furrowA)
	clk.Add(24*time.Hour - time.Second)
	if rec := send(furrowlight, "wk03-fri", furrowA); rec.Header().Get("Idempotent-Replayed") != "true" {
		t.Errorf("a retry 23h59m59s later = %d %s with no replay, want the kept answer", rec.Code, text(rec))
	}
	clk.Add(2 * time.Second)
	if rec := send(furrowlight, "wk03-fri", furrowB); rec.Code != http.StatusCreated {
		t.Errorf("the key with another body 24h0m1s later = %d %s, want 201: the key is new again", rec.Code, text(rec))
	}
}

func TestKeysBelongToTheirCustomer(t *testing.T) {
	as := bookings(newKeyStore((&clock{now: friday}).Now), NewMemStore())
	if rec := send(as("brinewatch"), "k1", brine3); rec.Code != http.StatusCreated {
		t.Fatalf("Brinewatch's POST /bookings with key k1 = %d %s, want 201", rec.Code, text(rec))
	}
	if rec := send(as("cirrine"), "k1", cirrine1); rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"bk-4502"`) {
		t.Errorf("Cirrine's POST /bookings with key k1 = %d %s, want 201 and bk-4502: Brinewatch's k1 is not Cirrine's", rec.Code, text(rec))
	}
}

// failOnce is a store whose first Create fails, the way a full disk does.
type failOnce struct {
	Store
	failed atomic.Bool
}

// Create fails the first time it is called, then keeps b.
func (s *failOnce) Create(ctx context.Context, b Booking) (Booking, error) {
	if s.failed.CompareAndSwap(false, true) {
		return Booking{}, errors.New("bookings.json: no space left on device")
	}
	return s.Store.Create(ctx, b)
}

func TestServerErrorIsNotKept(t *testing.T) {
	furrowlight := bookings(newKeyStore((&clock{now: friday}).Now), &failOnce{Store: NewMemStore()})("furrowlight")
	if rec := send(furrowlight, "wk03-fri", furrowA); rec.Code != http.StatusInternalServerError {
		t.Fatalf("POST /bookings with a full disk = %d %s, want 500", rec.Code, text(rec))
	}
	if rec := send(furrowlight, "wk03-fri", furrowA); rec.Code != http.StatusCreated {
		t.Errorf("the retry after a 500 = %d %s, want 201: a 500 is not kept, so the retry runs", rec.Code, text(rec))
	}
}

func TestNoKey(t *testing.T) {
	furrowlight := bookings(newKeyStore((&clock{now: friday}).Now), NewMemStore())("furrowlight")
	if rec := send(furrowlight, "", furrowA); rec.Code != http.StatusCreated {
		t.Fatalf("POST /bookings without a key = %d %s, want 201", rec.Code, text(rec))
	}
	if rec := send(furrowlight, "", furrowA); rec.Code != http.StatusConflict {
		t.Errorf("the same booking again without a key = %d %s, want 409 from the store, as before keys", rec.Code, text(rec))
	}
}
