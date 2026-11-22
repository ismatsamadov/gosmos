package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// elsewhere is a page on another site, like the one Mirela made on her laptop.
const elsewhere = "https://elsewhere.example"

// brinewatchKey is Brinewatch's test API key.
const brinewatchKey = "brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0"

// friday holds Furrowlight's booking of Monday's pass over Salar Quieto.
var friday = []Booking{
	{ID: "bk-4776", SatID: "furrow-b", Station: "SQT", AOS: time.Date(2028, 1, 24, 14, 19, 53, 0, time.UTC), LOS: time.Date(2028, 1, 24, 14, 29, 53, 0, time.UTC)},
}

// account returns the account page on Friday afternoon, 17:00 UTC, with Hollis signed in, and his session's ID.
func account(t *testing.T) (*Account, string) {
	t.Helper()
	store := NewMemStore()
	if err := store.load(friday); err != nil {
		t.Fatal(err)
	}
	now := func() time.Time { return time.Date(2028, 1, 21, 17, 0, 0, 0, time.UTC) }
	a := &Account{
		Store: store,
		Customers: map[string]Customer{
			"furrowlight": {ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}},
			"brinewatch":  {ID: "brinewatch", Name: "Brinewatch", SatIDs: []string{"brine-3"}, KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
		},
		Sessions: newSessions(now),
	}
	return a, a.Sessions.Start("furrowlight")
}

// browse sends one request through the guarded account page as a browser would: fetchSite and origin are the
// Sec-Fetch-Site and Origin headers it adds, left out when empty.
func browse(a *Account, method, path, fetchSite, origin, session, form string) *http.Response {
	r := httptest.NewRequest(method, path, strings.NewReader(form))
	if form != "" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if fetchSite != "" {
		r.Header.Set("Sec-Fetch-Site", fetchSite)
	}
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	if session != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: session})
	}
	rec := httptest.NewRecorder()
	noCrossOrigin(a.routes()).ServeHTTP(rec, r)
	return rec.Result()
}

func TestCancelFormFrom(t *testing.T) {
	tests := []struct {
		name              string
		fetchSite, origin string
		want              int
	}{
		{"Mirela's page on another site", "cross-site", elsewhere, http.StatusForbidden},
		{"a page elsewhere on the same site", "same-site", "https://news.example.com", http.StatusForbidden},
		{"an older browser, from another site", "", elsewhere, http.StatusForbidden},
		{"the account page's own form", "same-origin", "http://example.com", http.StatusSeeOther},
		{"an older browser, from the page itself", "", "http://example.com", http.StatusSeeOther},
		{"a program, which sends neither header", "", "", http.StatusSeeOther},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, hollis := account(t)
			res := browse(a, http.MethodPost, "/account/bookings/bk-4776/cancel", tt.fetchSite, tt.origin, hollis, "")
			if res.StatusCode != tt.want {
				t.Errorf("POST /account/bookings/bk-4776/cancel from %s = %d, want %d", tt.name, res.StatusCode, tt.want)
			}
			_, err := a.Store.Get(context.Background(), "bk-4776")
			if kept := err == nil; kept != (tt.want == http.StatusForbidden) {
				t.Errorf("after the form from %s, bk-4776 kept = %v, want %v", tt.name, kept, !kept)
			}
		})
	}
}

func TestALinkFromElsewhereStillOpensThePage(t *testing.T) {
	a, hollis := account(t)
	res := browse(a, http.MethodGet, "/account/", "cross-site", "", hollis, "")
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || !strings.Contains(string(body), "bk-4776") {
		t.Errorf("GET /account/ from a link on another site = %d, want 200 with bk-4776: reading the page changes nothing", res.StatusCode)
	}
}

func TestSignInFormFromElsewhere(t *testing.T) {
	a, _ := account(t)
	form := url.Values{"key": {brinewatchKey}}.Encode()
	res := browse(a, http.MethodPost, "/account/sign-in", "cross-site", elsewhere, "", form)
	if res.StatusCode != http.StatusForbidden || len(res.Cookies()) > 0 {
		t.Errorf("POST /account/sign-in from %s = %d with cookies %v, want 403 and none: no other site may sign a browser in", elsewhere, res.StatusCode, res.Cookies())
	}
	res = browse(a, http.MethodPost, "/account/sign-in", "same-origin", "http://example.com", "", form)
	if res.StatusCode != http.StatusSeeOther || len(res.Cookies()) != 1 {
		t.Errorf("POST /account/sign-in from the page's own form = %d with cookies %v, want 303 and a session cookie", res.StatusCode, res.Cookies())
	}
}
