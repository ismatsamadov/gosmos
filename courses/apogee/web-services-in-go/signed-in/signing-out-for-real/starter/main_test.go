package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// friday holds a booking of each of the two customers the tests sign in as.
var friday = []Booking{
	{ID: "bk-4776", SatID: "furrow-b", Station: "SQT", AOS: time.Date(2028, 1, 24, 14, 19, 53, 0, time.UTC), LOS: time.Date(2028, 1, 24, 14, 29, 53, 0, time.UTC)},
	{ID: "bk-4784", SatID: "pipit-1", Station: "PEL", AOS: time.Date(2028, 1, 24, 15, 8, 30, 0, time.UTC), LOS: time.Date(2028, 1, 24, 15, 17, 50, 0, time.UTC)},
}

// account returns the account page on Friday afternoon, 16:40 UTC.
func account(t *testing.T) *Account {
	t.Helper()
	store := NewMemStore()
	if err := store.load(friday); err != nil {
		t.Fatal(err)
	}
	now := func() time.Time { return time.Date(2028, 1, 21, 16, 40, 0, 0, time.UTC) }
	return &Account{
		Store: store,
		Customers: map[string]Customer{
			"apogee":      {ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}},
			"furrowlight": {ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}},
		},
		Sessions: newSessions(now),
	}
}

// send calls h with a request that carries this session cookie, if any, and returns the answer.
func send(t *testing.T, h http.HandlerFunc, method, path, session string) *http.Response {
	t.Helper()
	r := httptest.NewRequest(method, path, nil)
	if session != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: session})
	}
	rec := httptest.NewRecorder()
	func() {
		defer func() {
			if v := recover(); v != nil {
				t.Fatalf("%s %s panicked: %v", method, path, v)
			}
		}()
		h(rec, r)
	}()
	return rec.Result()
}

func TestSignOutForgetsTheSession(t *testing.T) {
	a := account(t)
	mirela := a.Sessions.Start("apogee")
	res := send(t, a.signOut, http.MethodPost, "/account/sign-out", mirela)
	if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/account/sign-in" {
		t.Fatalf("POST /account/sign-out = %d to %q, want 303 to %q", res.StatusCode, res.Header.Get("Location"), "/account/sign-in")
	}
	if customer, ok := a.Sessions.Lookup(mirela); ok {
		t.Errorf("after signing out, Lookup(Mirela's ID) = %q, %v, want \"\", false", customer, ok)
	}
	if res := send(t, a.page, http.MethodGet, "/account/", mirela); res.StatusCode != http.StatusSeeOther {
		t.Errorf("GET /account/ with Mirela's old cookie pasted back = %d, want 303: the service must forget the session", res.StatusCode)
	}
}

func TestSignOutClearsTheCookie(t *testing.T) {
	a := account(t)
	mirela := a.Sessions.Start("apogee")
	res := send(t, a.signOut, http.MethodPost, "/account/sign-out", mirela)
	var gone *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == "session" {
			gone = c
		}
	}
	if gone == nil {
		t.Fatalf("POST /account/sign-out sets cookies %v, want a session cookie that tells the browser to drop it", res.Cookies())
	}
	if gone.Path != "/account" {
		t.Errorf("sign-out's session cookie has Path %q, want %q: a cookie with another path is another cookie, and the browser keeps the old one", gone.Path, "/account")
	}
	if gone.MaxAge >= 0 || gone.Value != "" {
		t.Errorf("sign-out's Set-Cookie = %q, want an empty value with Max-Age=0", gone.Raw)
	}
}

func TestOtherBrowsersStaySignedIn(t *testing.T) {
	a := account(t)
	mirela := a.Sessions.Start("apogee")
	yusra := a.Sessions.Start("apogee")
	hollis := a.Sessions.Start("furrowlight")
	send(t, a.signOut, http.MethodPost, "/account/sign-out", mirela)
	for _, other := range []struct{ name, id string }{{"Yusra", yusra}, {"Hollis", hollis}} {
		if _, ok := a.Sessions.Lookup(other.id); !ok {
			t.Errorf("after Mirela signed out, Lookup(%s's ID) = false, want true: only her browser signed out", other.name)
		}
	}
}

func TestSignOutWithoutACookie(t *testing.T) {
	a := account(t)
	res := send(t, a.signOut, http.MethodPost, "/account/sign-out", "")
	if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/account/sign-in" {
		t.Errorf("POST /account/sign-out with no cookie = %d to %q, want 303 to %q", res.StatusCode, res.Header.Get("Location"), "/account/sign-in")
	}
}
