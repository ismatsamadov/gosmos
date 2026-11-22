package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// The booking service's test keys, readable on purpose: no production key ever appears in a lesson.
const (
	furrowKey = APIKey("furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b")
	brineKey  = APIKey("brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0")
)

// hexOf returns the 32 hex digits after a key's dot: the part that must never be logged.
func hexOf(k APIKey) string {
	return string(k)[strings.LastIndex(string(k), ".")+1:]
}

// logValue asks k for its LogValue, and stops the test when APIKey has none.
func logValue(t *testing.T, k APIKey) slog.Value {
	t.Helper()
	v, ok := any(k).(slog.LogValuer)
	if !ok {
		t.Fatalf("APIKey has no LogValue method, want one: func (k APIKey) LogValue() slog.Value")
	}
	return v.LogValue()
}

func TestLogValue(t *testing.T) {
	tests := []struct {
		name string
		k    APIKey
		want string
	}{
		{"Furrowlight's key", furrowKey, "furrowlight.[redacted]"},
		{"Brinewatch's key", brineKey, "brinewatch.[redacted]"},
		{"a key with no dot", APIKey(hexOf(furrowKey)), "[redacted]"},
		{"an empty key", APIKey(""), "[redacted]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := logValue(t, tt.k)
			if v.Kind() != slog.KindString || v.String() != tt.want {
				t.Errorf("APIKey(%q).LogValue() = %v, want the string %q", tt.k, v, tt.want)
			}
		})
	}
}

// TestKeyInEveryLine logs Furrowlight's key every way Mirela can think of and reads each line
// as anyone could: the secret must be nowhere, and the customer must still be named.
func TestKeyInEveryLine(t *testing.T) {
	o := &slog.HandlerOptions{ReplaceAttr: dropTime}
	tests := []struct {
		name  string
		write func(w *bytes.Buffer)
		want  string // the key as the line must show it
	}{
		{"JSON, as an attribute", func(w *bytes.Buffer) {
			slog.New(slog.NewJSONHandler(w, o)).Error("request failed", "key", furrowKey)
		}, `"key":"furrowlight.[redacted]"`},
		{"text, for people at a terminal", func(w *bytes.Buffer) {
			slog.New(slog.NewTextHandler(w, o)).Error("request failed", "key", furrowKey)
		}, `key=furrowlight.[redacted]`},
		{"JSON, inside a group", func(w *bytes.Buffer) {
			slog.New(slog.NewJSONHandler(w, o)).Error("request failed",
				slog.Group("customer", "id", "furrowlight", "key", furrowKey))
		}, `"customer":{"id":"furrowlight","key":"furrowlight.[redacted]"}`},
		{"JSON, through With", func(w *bytes.Buffer) {
			slog.New(slog.NewJSONHandler(w, o)).With("key", furrowKey).Error("request failed")
		}, `"key":"furrowlight.[redacted]"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.write(&buf)
			line := strings.TrimSuffix(buf.String(), "\n")
			if strings.Contains(line, hexOf(furrowKey)) {
				t.Errorf("the line holds the key's 32 hex digits (%s…), want them redacted: %s", hexOf(furrowKey)[:8], line)
			}
			if !strings.Contains(line, tt.want) {
				t.Errorf("the line = %s, want it to hold %s", line, tt.want)
			}
		})
	}
}

// TestKeyInJSONDecodes checks that the JSON line still decodes, the key a plain string in it.
func TestKeyInJSONDecodes(t *testing.T) {
	var buf bytes.Buffer
	slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime})).Error("request failed", "key", brineKey)
	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("the line %q isn't a JSON object (%v)", buf.String(), err)
	}
	if rec["key"] != "brinewatch.[redacted]" {
		b, _ := json.Marshal(rec["key"])
		t.Errorf(`"key" = %s, want "brinewatch.[redacted]"`, b)
	}
}
