package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// clock is a fake clock the tests move by hand.
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

// counted returns a handler that answers 200 and counts how many times it ran.
func counted(runs *int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*runs++
		w.Write([]byte("[]\n"))
	})
}

// get sends GET /bookings to h, with c in the request's context as authenticate leaves it, or with no customer when c is nil.
func get(h http.Handler, c *Customer) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, "/bookings", nil)
	if c != nil {
		r = r.WithContext(withCustomer(r.Context(), *c))
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

// headers writes down an answer's status and its three rate-limit headers.
func headers(rec *httptest.ResponseRecorder) string {
	h := rec.Header()
	return fmt.Sprintf("%d, X-RateLimit-Limit %q, X-RateLimit-Remaining %q, Retry-After %q",
		rec.Code, h.Get("X-RateLimit-Limit"), h.Get("X-RateLimit-Remaining"), h.Get("Retry-After"))
}

// errorBody returns the ErrorBody an answer holds, or what it holds instead.
func errorBody(rec *httptest.ResponseRecorder) string {
	var body ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		return fmt.Sprintf("%q", rec.Body)
	}
	return fmt.Sprintf("%+v", body)
}

func TestLimit(t *testing.T) {
	c := &clock{now: time.Date(2028, 1, 11, 9, 0, 0, 0, time.UTC)}
	runs := 0
	h := limit(newLimiter(c.Now))(counted(&runs))
	brinewatch := &Customer{ID: "brinewatch", Name: "Brinewatch"}
	for i := 1; i <= 10; i++ {
		want := fmt.Sprintf(`200, X-RateLimit-Limit "10", X-RateLimit-Remaining "%d", Retry-After ""`, 10-i)
		if got := headers(get(h, brinewatch)); got != want {
			t.Errorf("Brinewatch's request %d at 09:00:00 = %s, want %s", i, got, want)
		}
	}
	if runs != 10 {
		t.Errorf("the handler behind limit ran %d times for Brinewatch's first 10, want 10", runs)
	}
	steps := []struct {
		name  string
		after time.Duration
		want  string
	}{
		{"its 11th at 09:00:00", 0, `429, X-RateLimit-Limit "10", X-RateLimit-Remaining "0", Retry-After "1"`},
		{"0.6 of a token later", 600 * time.Millisecond, `429, X-RateLimit-Limit "10", X-RateLimit-Remaining "0", Retry-After "1"`},
		{"a second after that", time.Second, `200, X-RateLimit-Limit "10", X-RateLimit-Remaining "0", Retry-After ""`},
	}
	for _, step := range steps {
		c.Add(step.after)
		runs = 0
		rec := get(h, brinewatch)
		if got := headers(rec); got != step.want {
			t.Errorf("Brinewatch's request, %s = %s, want %s", step.name, got, step.want)
		}
		refused := rec.Code == http.StatusTooManyRequests
		if refused && runs != 0 {
			t.Errorf("Brinewatch's request, %s: the handler behind limit ran after a 429", step.name)
		}
		if want := (ErrorBody{Code: "rate_limited", Message: "slow down: 10 requests at once, then 1 a second"}); refused && errorBody(rec) != fmt.Sprintf("%+v", want) {
			t.Errorf("Brinewatch's request, %s: body = %s, want %+v", step.name, errorBody(rec), want)
		}
	}
	cirrine := &Customer{ID: "cirrine", Name: "Cirrine Atmospherics"}
	if got, want := headers(get(h, cirrine)), `200, X-RateLimit-Limit "10", X-RateLimit-Remaining "9", Retry-After ""`; got != want {
		t.Errorf("Cirrine's first request, behind Brinewatch's = %s, want %s", got, want)
	}
}

func TestLimitWithoutACustomer(t *testing.T) {
	c := &clock{now: time.Date(2028, 1, 11, 9, 0, 0, 0, time.UTC)}
	runs := 0
	rec := get(limit(newLimiter(c.Now))(counted(&runs)), nil)
	want := fmt.Sprintf("%+v", ErrorBody{Code: "internal", Message: "internal error"})
	if rec.Code != http.StatusInternalServerError || errorBody(rec) != want {
		t.Errorf("a request with no customer in its context = %d %s, want 500 %s", rec.Code, errorBody(rec), want)
	}
	if runs != 0 {
		t.Errorf("a request with no customer in its context: the handler behind limit ran %d times, want 0", runs)
	}
}
