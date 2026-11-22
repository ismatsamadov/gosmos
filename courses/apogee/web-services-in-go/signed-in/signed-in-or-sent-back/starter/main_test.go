package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

// friday holds the bookings the store keeps on Friday afternoon; bk-4784 is Yusra's test booking of Pipit-1.
var friday = []Booking{
	{ID: "bk-4776", SatID: "furrow-b", Station: "SQT", AOS: time.Date(2028, 1, 24, 14, 19, 53, 0, time.UTC), LOS: time.Date(2028, 1, 24, 14, 29, 53, 0, time.UTC)},
	{ID: "bk-4781", SatID: "cirrine-1", Station: "YST", AOS: time.Date(2028, 1, 21, 8, 51, 36, 0, time.UTC), LOS: time.Date(2028, 1, 21, 9, 1, 56, 0, time.UTC)},
	{ID: "bk-4783", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2028, 1, 21, 14, 31, 40, 0, time.UTC), LOS: time.Date(2028, 1, 21, 14, 41, 39, 0, time.UTC)},
	{ID: "bk-4784", SatID: "pipit-1", Station: "PEL", AOS: time.Date(2028, 1, 24, 15, 8, 30, 0, time.UTC), LOS: time.Date(2028, 1, 24, 15, 17, 50, 0, time.UTC)},
}

// account returns the account page as it stands on Friday afternoon, on the clock clk.
func account(t *testing.T, clk *clock) *Account {
	t.Helper()
	store := NewMemStore()
	if err := store.load(friday); err != nil {
		t.Fatal(err)
	}
	return &Account{
		Store: store,
		Customers: map[string]Customer{
			"apogee":      {ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}},
			"furrowlight": {ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}},
			"cirrine":     {ID: "cirrine", Name: "Cirrine Atmospherics", SatIDs: []string{"cirrine-1"}},
		},
		Sessions: newSessions(clk.Now),
	}
}

// open asks for the account page with this Cookie header, as a browser does, and returns the answer and its body.
func open(t *testing.T, a *Account, cookie string) (*http.Response, string) {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, "/account/", nil)
	if cookie != "" {
		r.Header.Set("Cookie", cookie)
	}
	rec := httptest.NewRecorder()
	func() {
		defer func() {
			if v := recover(); v != nil {
				t.Fatalf("GET /account/ with Cookie %q panicked: %v", cookie, v)
			}
		}()
		a.page(rec, r)
	}()
	res := rec.Result()
	body, _ := io.ReadAll(res.Body)
	return res, string(body)
}

func TestSignedIn(t *testing.T) {
	a := account(t, &clock{now: time.Date(2028, 1, 21, 16, 10, 0, 0, time.UTC)})
	yusra := a.Sessions.Start("apogee")
	stranger := a.Sessions.Start("brinewatch") // a customer this page doesn't know
	tests := []struct {
		name, cookie string
		in           bool
	}{
		{"Yusra's session", "session=" + yusra, true},
		{"Yusra's session after another cookie", "units=km; session=" + yusra, true},
		{"no cookie at all", "", false},
		{"only another cookie", "units=km", false},
		{"a session nobody started", "session=ZQ7XB2MZ6KD3WPLT5HNC7YFEJ4", false},
		{"Apogee's ID typed as the session", "session=apogee", false},
		{"an empty session cookie", "session=", false},
		{"the session of a customer the page doesn't know", "session=" + stranger, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, body := open(t, a, tt.cookie)
			if !tt.in {
				if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/account/sign-in" {
					t.Errorf("GET /account/ with Cookie %q = %d to %q, want 303 to %q", tt.cookie, res.StatusCode, res.Header.Get("Location"), "/account/sign-in")
				}
				return
			}
			if res.StatusCode != http.StatusOK || !strings.Contains(body, "bk-4784") {
				t.Fatalf("GET /account/ with Cookie %q = %d, want 200 and Yusra's bk-4784:\n%s", tt.cookie, res.StatusCode, body)
			}
			for _, other := range []string{"bk-4776", "bk-4781", "bk-4783"} {
				if strings.Contains(body, other) {
					t.Errorf("GET /account/ as Apogee shows %s, another customer's booking", other)
				}
			}
			renewed := ""
			for _, c := range res.Cookies() {
				if c.Name == "session" {
					renewed = c.Value
				}
			}
			if renewed != yusra {
				t.Errorf("GET /account/ renews the cookie as session=%s, want session=%s: signedIn returns the cookie's value first", renewed, yusra)
			}
		})
	}
}

func TestEndedSessionSendsBack(t *testing.T) {
	clk := &clock{now: time.Date(2028, 1, 21, 16, 10, 0, 0, time.UTC)}
	a := account(t, clk)
	hollis := a.Sessions.Start("furrowlight")
	res, body := open(t, a, "session="+hollis)
	if res.StatusCode != http.StatusOK || !strings.Contains(body, "bk-4776") || !strings.Contains(body, "bk-4783") {
		t.Fatalf("GET /account/ as Furrowlight = %d, want 200 with bk-4776 and bk-4783:\n%s", res.StatusCode, body)
	}
	clk.Add(30 * time.Minute)
	res, _ = open(t, a, "session="+hollis)
	if res.StatusCode != http.StatusSeeOther {
		t.Errorf("GET /account/ 30 minutes after Hollis's last visit = %d, want 303 to sign in again", res.StatusCode)
	}
}
