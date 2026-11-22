package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// The test secrets of the booking service and the ops console, readable on purpose.
const (
	furrowKey   = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	furrowToken = "furrowlight.1830682800.767de44da59ef9b9e3afb9354270dd7cc56e26ceb00650b8a615afe1cd10b616"
	password    = "fails closed, always"
	hookSecret  = "furrowlight webhook secret"
)

func TestRedact(t *testing.T) {
	at := time.Date(2028, 11, 27, 10, 15, 0, 0, time.UTC)
	tests := []struct {
		name   string
		groups []string
		a      slog.Attr
		blank  bool // whether redact must replace the value
	}{
		{"an API key at the top level", nil, slog.String("api_key", furrowKey), true},
		{"a header dumped as attributes", []string{"header"}, slog.String("Authorization", "Bearer "+furrowToken), true},
		{"a password in capitals", []string{"console", "login"}, slog.String("PASSWORD", password), true},
		{"a webhook's secret", []string{"webhook"}, slog.String("secret", hookSecret), true},
		{"a token two groups down", []string{"auth", "bearer"}, slog.String("Token", furrowToken), true},
		{"a count of tokens", nil, slog.Int("token_count", 3), false},
		{"a pass ID", nil, slog.String("pass_id", "cirrine-1/LAN/16402"), false},
		{"the built-in time", nil, slog.Time(slog.TimeKey, at), false},
		{"the built-in level", nil, slog.Any(slog.LevelKey, slog.LevelInfo), false},
		{"the built-in message", nil, slog.String(slog.MessageKey, "request failed"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := redact(tt.groups, tt.a)
			want := tt.a
			if tt.blank {
				want = slog.String(tt.a.Key, "[redacted]")
			}
			if !got.Equal(want) {
				t.Errorf("redact(%q, %v) = %v, want %v", tt.groups, tt.a, got, want)
			}
		})
	}
}

// TestRedactInHandler writes one of the records Mirela found, its time fixed, through a JSON
// handler with redact, and reads the line as anyone could.
func TestRedactInHandler(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: redact})
	r := slog.NewRecord(time.Date(2028, 11, 22, 14, 3, 7, 0, time.UTC), slog.LevelDebug, "request headers", 0)
	r.AddAttrs(
		slog.String("api_key", furrowKey),
		slog.Group("header",
			slog.String("Authorization", "Bearer "+furrowToken),
			slog.String("Content-Type", "application/json")),
		slog.Group("console", slog.Group("login",
			slog.String("user", "mirela"),
			slog.String("Password", password))),
		slog.Int("token_count", 3),
	)
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	line := strings.TrimSuffix(buf.String(), "\n")
	for _, sc := range []struct{ what, s string }{
		{"Furrowlight's key", furrowKey},
		{"Furrowlight's token", furrowToken},
		{"the password", password},
	} {
		if strings.Contains(line, sc.s) {
			t.Errorf("the line holds %s, want it redacted: %s", sc.what, line)
		}
	}
	var rec map[string]any
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		t.Fatalf("the line %q isn't a JSON object (%v)", line, err)
	}
	for _, f := range []struct {
		path []string
		want any
	}{
		{[]string{"time"}, "2028-11-22T14:03:07Z"},
		{[]string{"level"}, "DEBUG"},
		{[]string{"msg"}, "request headers"},
		{[]string{"api_key"}, "[redacted]"},
		{[]string{"header", "Authorization"}, "[redacted]"},
		{[]string{"header", "Content-Type"}, "application/json"},
		{[]string{"console", "login", "user"}, "mirela"},
		{[]string{"console", "login", "Password"}, "[redacted]"},
		{[]string{"token_count"}, 3.0},
	} {
		var v any = rec
		for _, k := range f.path {
			m, _ := v.(map[string]any)
			v = m[k]
		}
		if v != f.want {
			b, _ := json.Marshal(v)
			w, _ := json.Marshal(f.want)
			t.Errorf("%q = %s, want %s", strings.Join(f.path, "."), b, w)
		}
	}
}
