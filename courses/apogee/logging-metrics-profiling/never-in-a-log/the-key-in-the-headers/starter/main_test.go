package main

import (
	"bytes"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"strings"
	"testing"
)

// The booking service's test fixtures: Furrowlight's API key, its token that expires at 11:00
// on 5 January 2028, and the signature of its launch-day event under its test webhook secret.
const (
	furrowKey   = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	furrowToken = "furrowlight.1830682800.767de44da59ef9b9e3afb9354270dd7cc56e26ceb00650b8a615afe1cd10b616"
	eventSig    = "367bc6bb2684c3d410a5ab09b86be50e0e5b66ca8740205ded7f6bebffe4314f"
)

// secrets are the hex digits of each fixture: none may appear in a logged line.
var secrets = []struct{ what, digits string }{
	{"the key's 32 hex digits", "5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"},
	{"the token's 64 hex digits", "767de44da59ef9b9e3afb9354270dd7cc56e26ceb00650b8a615afe1cd10b616"},
	{"the signature's 64 hex digits", eventSig},
}

// logged writes h as bookingd's DEBUG line does, through safeHeader, and returns the line,
// with no time, as anyone reading logs in the ops console would see it.
func logged(h http.Header) string {
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: dropTime}))
	l.Debug("request headers", "header", safeHeader(h))
	return buf.String()
}

func TestSafeHeaderRedacts(t *testing.T) {
	tests := []struct {
		name string
		h    func() http.Header
		want map[string][]string // the copy's values, name by name
	}{
		{"Furrowlight asks for a token", func() http.Header {
			h := http.Header{}
			h.Set("Authorization", "Key "+furrowKey)
			h.Set("Content-Type", "application/json")
			return h
		}, map[string][]string{"Authorization": {"[redacted]"}, "Content-Type": {"application/json"}}},
		{"Furrowlight books with its token", func() http.Header {
			h := http.Header{}
			h.Set("Authorization", "Bearer "+furrowToken)
			h.Set("Idempotency-Key", "wk49-mon")
			return h
		}, map[string][]string{"Authorization": {"[redacted]"}, "Idempotency-Key": {"wk49-mon"}}},
		{"a client that sent its key and its token", func() http.Header {
			h := http.Header{}
			h.Add("Authorization", "Key "+furrowKey)
			h.Add("Authorization", "Bearer "+furrowToken)
			return h
		}, map[string][]string{"Authorization": {"[redacted]", "[redacted]"}}},
		{"a webhook's signature", func() http.Header {
			h := http.Header{}
			h.Set("Apogee-Timestamp", "1833027183")
			h.Set("Apogee-Signature", eventSig)
			return h
		}, map[string][]string{"Apogee-Signature": {"[redacted]"}, "Apogee-Timestamp": {"1833027183"}}},
		{"a token in a cookie, added in lowercase", func() http.Header {
			h := http.Header{}
			h.Add("cookie", "token="+furrowToken)
			h.Add("accept", "application/json")
			return h
		}, map[string][]string{"Cookie": {"[redacted]"}, "Accept": {"application/json"}}},
		{"an answer that sets two cookies", func() http.Header {
			h := http.Header{}
			h.Add("Set-Cookie", "token="+furrowToken+"; HttpOnly")
			h.Add("Set-Cookie", "guide=1")
			return h
		}, map[string][]string{"Set-Cookie": {"[redacted]", "[redacted]"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := safeHeader(tt.h())
			for _, name := range slices.Sorted(maps.Keys(tt.want)) {
				if want := tt.want[name]; !slices.Equal(got[name], want) {
					t.Errorf("safeHeader(h)[%q] = %q, want %q", name, got[name], want)
				}
			}
			if len(got) != len(tt.want) {
				t.Errorf("safeHeader(h) has %d headers, want %d: keep every name, add none", len(got), len(tt.want))
			}
			if t.Failed() {
				return
			}
			line := logged(tt.h())
			for _, sc := range secrets {
				if strings.Contains(line, sc.digits) {
					t.Errorf("the logged line holds %s (%s…), want them redacted: %s", sc.what, sc.digits[:8], line)
				}
			}
		})
	}
}

func TestSafeHeaderLeavesTheRequest(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Key "+furrowKey)
	h.Add("Cookie", "token="+furrowToken)
	h.Add("Cookie", "guide=1")
	safeHeader(h)
	if got, want := h.Values("Authorization"), []string{"Key " + furrowKey}; !slices.Equal(got, want) {
		t.Errorf(`after safeHeader(h), h["Authorization"] = %q, want %q: bookingd still checks the request's own key`, got, want)
	}
	if got, want := h.Values("Cookie"), []string{"token=" + furrowToken, "guide=1"}; !slices.Equal(got, want) {
		t.Errorf(`after safeHeader(h), h["Cookie"] = %q, want %q: change only the copy`, got, want)
	}
}
