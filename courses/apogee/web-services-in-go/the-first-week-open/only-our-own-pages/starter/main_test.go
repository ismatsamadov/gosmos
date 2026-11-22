package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGoStaysOnTheSite(t *testing.T) {
	var logs bytes.Buffer
	h := requestID(goNext(log.New(&logs, "", 0)))
	cases := []struct {
		name, next, want string
	}{
		{"the guide's page on files", "/guide/files.html", "/guide/files.html"},
		{"a page with a query and a fragment", "/guide/webhooks.html?from=launch#retries", "/guide/webhooks.html?from=launch#retries"},
		{"the account page", "/account/", "/account/"},
		{"no next at all", "", "/guide/"},
		{"Mirela's link", "https://apogee-bookings.example/renew", "/guide/"},
		{"another site without its scheme", "//apogee-bookings.example/renew", "/guide/"},
		{"three slashes", "///apogee-bookings.example/renew", "/guide/"},
		{"a backslash after the slash", `/\apogee-bookings.example/renew`, "/guide/"},
		{"two backslashes", `\\apogee-bookings.example/renew`, "/guide/"},
		{"a tab between the slashes", "/\t/apogee-bookings.example/renew", "/guide/"},
		{"a line break between the slashes", "/\n/apogee-bookings.example/renew", "/guide/"},
		{"a scheme and one slash", "https:/apogee-bookings.example/renew", "/guide/"},
		{"a script", "javascript:alert(1)", "/guide/"},
		{"a page without its first slash", "guide/files.html", "/guide/"},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodGet, "/go?"+url.Values{"next": {c.next}}.Encode(), nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		if got := rec.Header().Get("Location"); rec.Code != http.StatusSeeOther || got != c.want {
			t.Errorf("GET /go with next %q, %s = %d to %q, want 303 to %q", c.next, c.name, rec.Code, got, c.want)
		}
	}
	for _, line := range strings.Split(strings.TrimSpace(logs.String()), "\n") {
		if strings.Contains(line, "apogee-bookings.example") {
			t.Errorf("log line %q names Mirela's site, want only pages of this site in the log", line)
		}
	}
}
