package main

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// inJSON returns a as the JSON handler writes it inside a record, such as "wait_ms":570000.
func inJSON(a slog.Attr) string {
	bare := func(groups []string, x slog.Attr) slog.Attr {
		if len(groups) == 0 && (x.Key == slog.TimeKey || x.Key == slog.LevelKey || x.Key == slog.MessageKey) {
			return slog.Attr{}
		}
		return x
	}
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: bare}))
	l.LogAttrs(context.Background(), slog.LevelError, "pass failed", a)
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(buf.String()), "{"), "}")
}

func TestMsAttr(t *testing.T) {
	tests := []struct {
		name string
		key  string
		d    time.Duration
		want string
	}{
		{"the 01:55 pass's wait", "wait", 9*time.Minute + 30*time.Second, `"wait_ms":570000`},
		{"an acknowledgment after a second and a half", "ack", 1500 * time.Millisecond, `"ack_ms":1500`},
		{"999 microseconds, rounded down", "ack", 999 * time.Microsecond, `"ack_ms":0`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inJSON(msAttr(tt.key, tt.d)); got != tt.want {
				t.Errorf("msAttr(%q, %v) = %s, want %s", tt.key, tt.d, got, tt.want)
			}
		})
	}
}
