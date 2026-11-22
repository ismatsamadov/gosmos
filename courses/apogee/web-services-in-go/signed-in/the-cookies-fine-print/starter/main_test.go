package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// id is a session ID such as the service will make: 26 random characters.
const id = "J4RQ7XB2MZ6KD3WPLT5HNC7YFE"

// sameSite names a SameSite value as Go writes it.
func sameSite(v http.SameSite) string {
	switch v {
	case 0:
		return "0 (no SameSite at all)"
	case http.SameSiteDefaultMode:
		return "http.SameSiteDefaultMode"
	case http.SameSiteLaxMode:
		return "http.SameSiteLaxMode"
	case http.SameSiteStrictMode:
		return "http.SameSiteStrictMode"
	case http.SameSiteNoneMode:
		return "http.SameSiteNoneMode"
	}
	return "an unknown SameSite"
}

func TestSessionCookie(t *testing.T) {
	c := sessionCookie(id)
	if c.Name != "session" || c.Value != id {
		t.Fatalf("sessionCookie(id) = %s=%s, want session=%s", c.Name, c.Value, id)
	}
	if c.Path != "/account" {
		t.Errorf("sessionCookie(id).Path = %q, want %q: the API never needs the cookie", c.Path, "/account")
	}
	if c.MaxAge != 1800 {
		t.Errorf("sessionCookie(id).MaxAge = %v, want 1800: 30 minutes, counted in seconds", c.MaxAge)
	}
	if !c.HttpOnly {
		t.Errorf("sessionCookie(id).HttpOnly = %v, want true: no script may read the session", c.HttpOnly)
	}
	if !c.Secure {
		t.Errorf("sessionCookie(id).Secure = %v, want true: the cookie travels only over HTTPS", c.Secure)
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Errorf("sessionCookie(id).SameSite = %s, want http.SameSiteLaxMode", sameSite(c.SameSite))
	}
}

func TestTheHeaderTheBrowserGets(t *testing.T) {
	rec := httptest.NewRecorder()
	http.SetCookie(rec, sessionCookie(id))
	want := "session=" + id + "; Path=/account; Max-Age=1800; HttpOnly; Secure; SameSite=Lax"
	if got := rec.Header().Get("Set-Cookie"); got != want {
		t.Errorf("Set-Cookie = %q, want %q", got, want)
	}
}
