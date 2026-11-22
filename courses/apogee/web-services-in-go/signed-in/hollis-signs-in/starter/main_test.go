package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// The test API keys of the two customers who sign in tonight.
const (
	apogeeKey      = "apogee.9e4d2a7c1b0f3e5d8a6c4b2f0e1d3a5c"
	furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
)

// planningPage is the one origin the API trusts: Furrowlight's planning page.
const planningPage = "https://plan.furrowlight.example"

// elsewhere is a page on another site, like the one Mirela made to try the account page.
const elsewhere = "https://elsewhere.example"

// token is the token Furrowlight's server got from POST /tokens at 18:00 UTC and handed to the planning page.
const token = "furrowlight.1832094000.86f002a347c6c43044455ec26a704ffdfb2fa2cd06751c50744abfa8422fde1e"

// customers returns the customer table as the service keeps it: a hash of each key, never the key.
func customers() map[string]Customer {
	return map[string]Customer{
		"apogee":      {ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}, KeyHash: "8018b533030c734641dda97af10c84dd57357d0e37b7a0bdee7ed2a63c911407"},
		"furrowlight": {ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}, KeyHash: "95ff036b9ad2e59c26fae1be7dda78b12412d243fcbb6567bc8124e6f1c93fd3"},
		"brinewatch":  {ID: "brinewatch", Name: "Brinewatch", SatIDs: []string{"brine-3"}, KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
		"cirrine":     {ID: "cirrine", Name: "Cirrine Atmospherics", SatIDs: []string{"cirrine-1"}, KeyHash: "d9ee2e5de1796bbeba223c0737933b3e60d53687d479f8a04c98c99df87b510d"},
	}
}

// friday holds the bookings the store keeps on Friday evening.
var friday = []Booking{
	{ID: "bk-4776", SatID: "furrow-b", Station: "SQT", AOS: time.Date(2028, 1, 24, 14, 19, 53, 0, time.UTC), LOS: time.Date(2028, 1, 24, 14, 29, 53, 0, time.UTC)},
	{ID: "bk-4781", SatID: "cirrine-1", Station: "YST", AOS: time.Date(2028, 1, 21, 8, 51, 36, 0, time.UTC), LOS: time.Date(2028, 1, 21, 9, 1, 56, 0, time.UTC)},
	{ID: "bk-4783", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2028, 1, 21, 14, 31, 40, 0, time.UTC), LOS: time.Date(2028, 1, 21, 14, 41, 39, 0, time.UTC)},
	{ID: "bk-4784", SatID: "pipit-1", Station: "PEL", AOS: time.Date(2028, 1, 24, 15, 8, 30, 0, time.UTC), LOS: time.Date(2028, 1, 24, 15, 17, 50, 0, time.UTC)},
}

// api stands in for the booking service's routes: it answers a caller with Furrowlight's token and refuses any other.
func api() http.Handler {
	mux := http.NewServeMux()
	withToken := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+token {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, `{"error":"unauthorized","message":"missing credentials"}`, http.StatusUnauthorized)
				return
			}
			h(w, r)
		}
	}
	mux.HandleFunc("GET /passes", withToken(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "[]\n")
	}))
	mux.HandleFunc("POST /bookings", withToken(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/bookings/bk-4785")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"id":"bk-4785"}`+"\n")
	}))
	return mux
}

// fridayEvening returns the whole site as it opens to the beta, and the store behind it.
func fridayEvening(t *testing.T) (http.Handler, *MemStore) {
	t.Helper()
	store := NewMemStore()
	if err := store.load(friday); err != nil {
		t.Fatal(err)
	}
	clk := &clock{now: time.Date(2028, 1, 21, 18, 0, 0, 0, time.UTC)}
	a := &Account{Store: store, Customers: customers(), Sessions: newSessions(clk.Now)}
	return site(api(), a, planningPage), store
}

// visit is one request as a browser sends it: from a page at origin, which may be the service's own.
func visit(method, path, origin, session, form string) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(form))
	if form != "" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if origin == "" {
		r.Header.Set("Sec-Fetch-Site", "same-origin")
	} else {
		r.Header.Set("Sec-Fetch-Site", "cross-site")
		r.Header.Set("Origin", origin)
	}
	if session != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: session})
	}
	return r
}

// serve sends r through h and returns the answer as the browser gets it.
func serve(h http.Handler, r *http.Request) *http.Response {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec.Result()
}

// signIn signs a browser in with key from the page's own form and returns the new session's ID.
func signIn(t *testing.T, h http.Handler, key string) string {
	t.Helper()
	res := serve(h, visit(http.MethodPost, "/account/sign-in", "", "", url.Values{"key": {key}}.Encode()))
	for _, c := range res.Cookies() {
		if c.Name == "session" && c.Value != "" && res.StatusCode == http.StatusSeeOther {
			return c.Value
		}
	}
	t.Fatalf("POST /account/sign-in = %d with cookies %v, want 303 with a session cookie", res.StatusCode, res.Cookies())
	return ""
}

// body returns an answer's body as text.
func body(res *http.Response) string {
	b, _ := io.ReadAll(res.Body)
	return string(b)
}

// guarded checks an answer of the account page for Mirela's four headers and no-store.
func guarded(t *testing.T, call string, res *http.Response) {
	t.Helper()
	for _, h := range []struct{ name, want string }{
		{"Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'"},
		{"Strict-Transport-Security", "max-age=31536000"},
		{"X-Content-Type-Options", "nosniff"},
		{"Referrer-Policy", "same-origin"},
		{"Cache-Control", "no-store"},
	} {
		if got := res.Header.Get(h.name); got != h.want {
			t.Errorf("%s: %s = %q, want %q", call, h.name, got, h.want)
		}
	}
}

// kept reports whether the store still holds booking id.
func kept(store *MemStore, id string) bool {
	_, err := store.Get(context.Background(), id)
	return err == nil
}

func TestHollisSignsIn(t *testing.T) {
	h, _ := fridayEvening(t)
	res := serve(h, visit(http.MethodPost, "/account/sign-in", "", "", url.Values{"key": {furrowlightKey}}.Encode()))
	if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/account/" {
		t.Fatalf("POST /account/sign-in = %d to %q, want 303 to %q", res.StatusCode, res.Header.Get("Location"), "/account/")
	}
	guarded(t, "POST /account/sign-in", res)
	id := signIn(t, h, furrowlightKey)
	res = serve(h, visit(http.MethodGet, "/account/", "", id, ""))
	page := body(res)
	if res.StatusCode != http.StatusOK || !strings.Contains(page, "bk-4776") || !strings.Contains(page, "bk-4783") {
		t.Fatalf("GET /account/ as Furrowlight = %d, want 200 with bk-4776 and bk-4783:\n%s", res.StatusCode, page)
	}
	for _, other := range []string{"bk-4781", "bk-4784"} {
		if strings.Contains(page, other) {
			t.Errorf("GET /account/ as Furrowlight shows %s, which is another customer's", other)
		}
	}
	guarded(t, "GET /account/", res)
	res = serve(h, visit(http.MethodGet, "/account/nowhere", "", id, ""))
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("GET /account/nowhere = %d, want 404", res.StatusCode)
	}
	guarded(t, "GET /account/nowhere", res)
}

func TestFormsFromElsewhere(t *testing.T) {
	h, store := fridayEvening(t)
	hollis := signIn(t, h, furrowlightKey)
	res := serve(h, visit(http.MethodPost, "/account/bookings/bk-4776/cancel", elsewhere, hollis, ""))
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("POST /account/bookings/bk-4776/cancel from %s = %d, want 403", elsewhere, res.StatusCode)
	}
	if !kept(store, "bk-4776") {
		t.Errorf("after a form from %s, bk-4776 is gone, want it kept", elsewhere)
	}
	guarded(t, "the 403 for a form from "+elsewhere, res)
	yusra := signIn(t, h, apogeeKey)
	res = serve(h, visit(http.MethodPost, "/account/bookings/bk-4784/cancel", "", yusra, ""))
	if res.StatusCode != http.StatusSeeOther || kept(store, "bk-4784") {
		t.Errorf("POST /account/bookings/bk-4784/cancel from the page itself = %d, kept %v; want 303 and bk-4784 gone", res.StatusCode, kept(store, "bk-4784"))
	}
}

func TestPlanningPage(t *testing.T) {
	h, _ := fridayEvening(t)
	pre := visit(http.MethodOptions, "/passes", planningPage, "", "")
	pre.Header.Set("Access-Control-Request-Method", "GET")
	pre.Header.Set("Access-Control-Request-Headers", "authorization")
	res := serve(h, pre)
	if res.StatusCode != http.StatusNoContent || res.Header.Get("Access-Control-Allow-Origin") != planningPage {
		t.Fatalf("preflight OPTIONS /passes from %s = %d with Access-Control-Allow-Origin %q, want 204 with %q: cors answers it",
			planningPage, res.StatusCode, res.Header.Get("Access-Control-Allow-Origin"), planningPage)
	}
	get := visit(http.MethodGet, "/passes", planningPage, "", "")
	get.Header.Set("Authorization", "Bearer "+token)
	res = serve(h, get)
	if res.StatusCode != http.StatusOK || res.Header.Get("Access-Control-Allow-Origin") != planningPage {
		t.Errorf("GET /passes from %s = %d with Access-Control-Allow-Origin %q, want 200 with %q",
			planningPage, res.StatusCode, res.Header.Get("Access-Control-Allow-Origin"), planningPage)
	}
	post := visit(http.MethodPost, "/bookings", planningPage, "", "")
	post.Header.Set("Authorization", "Bearer "+token)
	post.Header.Set("Content-Type", "application/json")
	res = serve(h, post)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("POST /bookings from %s with a token = %d, want 201 from the API: a token is never sent by itself, so the API needs no form guard", planningPage, res.StatusCode)
	}
	page := visit(http.MethodGet, "/account/sign-in", planningPage, "", "")
	res = serve(h, page)
	if got := res.Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("GET /account/sign-in from %s: Access-Control-Allow-Origin = %q, want none: only the API answers other origins", planningPage, got)
	}
	res = serve(h, visit(http.MethodPut, "/passes", "", "", ""))
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("PUT /passes = %d, want 405 from the API's own mux", res.StatusCode)
	}
}

func TestSignOut(t *testing.T) {
	h, _ := fridayEvening(t)
	id := signIn(t, h, furrowlightKey)
	res := serve(h, visit(http.MethodPost, "/account/sign-out", "", id, ""))
	if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/account/sign-in" {
		t.Fatalf("POST /account/sign-out = %d to %q, want 303 to %q", res.StatusCode, res.Header.Get("Location"), "/account/sign-in")
	}
	res = serve(h, visit(http.MethodGet, "/account/", "", id, ""))
	if res.StatusCode != http.StatusSeeOther {
		t.Errorf("GET /account/ with the cookie from before signing out = %d, want 303 to sign in again", res.StatusCode)
	}
}
