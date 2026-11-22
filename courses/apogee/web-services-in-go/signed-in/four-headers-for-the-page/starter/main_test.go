package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"
)

// friday holds Yusra's test booking of Pipit-1 over Port Ellery on Monday.
var friday = []Booking{
	{ID: "bk-4784", SatID: "pipit-1", Station: "PEL", AOS: time.Date(2028, 1, 24, 15, 8, 30, 0, time.UTC), LOS: time.Date(2028, 1, 24, 15, 17, 50, 0, time.UTC)},
}

// account returns the account page on Friday afternoon, 17:40 UTC, with Yusra signed in as Apogee, and her session's ID.
func account(t *testing.T) (*Account, string) {
	t.Helper()
	store := NewMemStore()
	if err := store.load(friday); err != nil {
		t.Fatal(err)
	}
	now := func() time.Time { return time.Date(2028, 1, 21, 17, 40, 0, 0, time.UTC) }
	a := &Account{
		Store:     store,
		Customers: map[string]Customer{"apogee": {ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}}},
		Sessions:  newSessions(now),
	}
	return a, a.Sessions.Start("apogee")
}

// request returns a request from the page's own browser, with Yusra's session cookie when session isn't empty.
func request(method, path, session, form string) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(form))
	if form != "" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	if session != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: session})
	}
	return r
}

// policy returns the directives of a Content-Security-Policy, sorted, so their order doesn't matter.
func policy(csp string) []string {
	var directives []string
	for _, d := range strings.Split(csp, ";") {
		if d = strings.TrimSpace(d); d != "" {
			directives = append(directives, d)
		}
	}
	slices.Sort(directives)
	return directives
}

func TestEveryAnswerCarriesTheHeaders(t *testing.T) {
	a, yusra := account(t)
	fromElsewhere := request(http.MethodPost, "/account/bookings/bk-4784/cancel", yusra, "")
	fromElsewhere.Header.Set("Sec-Fetch-Site", "cross-site")
	fromElsewhere.Header.Set("Origin", "https://elsewhere.example")
	tests := []struct {
		name   string
		r      *http.Request
		status int
	}{
		{"the page, signed in", request(http.MethodGet, "/account/", yusra, ""), http.StatusOK},
		{"the page, signed out", request(http.MethodGet, "/account/", "", ""), http.StatusSeeOther},
		{"the sign-in form", request(http.MethodGet, "/account/sign-in", "", ""), http.StatusOK},
		{"a wrong key", request(http.MethodPost, "/account/sign-in", "", url.Values{"key": {"apogee.0000"}}.Encode()), http.StatusUnauthorized},
		{"a form from another site", fromElsewhere, http.StatusForbidden},
		{"an address the page doesn't have", request(http.MethodGet, "/account/nowhere", "", ""), http.StatusNotFound},
	}
	h := securityHeaders(noCrossOrigin(a.routes()))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, tt.r)
			res := rec.Result() // the headers as they were when the answer was written
			call := tt.r.Method + " " + tt.r.URL.Path
			if res.StatusCode != tt.status {
				t.Errorf("%s = %d, want %d", call, res.StatusCode, tt.status)
			}
			want := []string{"default-src 'self'", "frame-ancestors 'none'"}
			if got := res.Header.Get("Content-Security-Policy"); !slices.Equal(policy(got), want) {
				t.Errorf("%s: Content-Security-Policy = %q, want %q", call, got, "default-src 'self'; frame-ancestors 'none'")
			}
			for _, h := range []struct{ name, want string }{
				{"Strict-Transport-Security", "max-age=31536000"},
				{"X-Content-Type-Options", "nosniff"},
				{"Referrer-Policy", "same-origin"},
			} {
				if got := res.Header.Get(h.name); got != h.want {
					t.Errorf("%s: %s = %q, want %q", call, h.name, got, h.want)
				}
			}
		})
	}
}
