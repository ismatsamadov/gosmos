package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// cirrine is the pass ID of Cirrine-1's pass of 8 November.
const cirrine = "cirrine-1/LAN/16402"

// logger returns a logger on passHandler around a JSON handler at INFO that writes to buf, without the time.
func logger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(passHandler{slog.NewJSONHandler(buf, &slog.HandlerOptions{ReplaceAttr: dropTime})})
}

// booking is the context of bookingd's request req-004211, which booked Cirrine-1's pass.
func booking() context.Context {
	return context.WithValue(withPass(context.Background(), cirrine), requestIDKey{}, "req-004211")
}

// decode reads the one line in buf; it stops the test, naming call, unless there is exactly one.
func decode(t *testing.T, call string, buf *bytes.Buffer) (map[string]any, string) {
	t.Helper()
	raw := strings.TrimSuffix(buf.String(), "\n")
	var rec map[string]any
	if raw == "" || strings.Contains(raw, "\n") || json.Unmarshal([]byte(raw), &rec) != nil {
		t.Fatalf("%s wrote %q, want one JSON line", call, buf.String())
	}
	return rec, raw
}

// field writes the value of rec's key as JSON writes it, or says it is missing.
func field(rec map[string]any, key string) string {
	v, ok := rec[key]
	if !ok {
		return "missing"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func TestHandle(t *testing.T) {
	t.Run("both IDs, after the record's own", func(t *testing.T) {
		var buf bytes.Buffer
		const call = `InfoContext(ctx, "booking created", "booking_id", "bk-5611")`
		logger(&buf).InfoContext(booking(), "booking created", "booking_id", "bk-5611")
		rec, raw := decode(t, call, &buf)
		if rec["pass_id"] != cirrine || rec["request_id"] != "req-004211" {
			t.Fatalf(`%s: "pass_id" = %s and "request_id" = %s, want %q and "req-004211"`, call, field(rec, "pass_id"), field(rec, "request_id"), cirrine)
		}
		if own, p := strings.Index(raw, `"booking_id":`), strings.Index(raw, `"pass_id":`); p < own {
			t.Errorf("%s wrote %s, want pass_id after the record's own booking_id", call, raw)
		}
		if n := strings.Count(raw, `"pass_id":`); n != 1 {
			t.Errorf("%s wrote %d pass_id keys, want 1: %s", call, n, raw)
		}
	})
	t.Run("only the pass", func(t *testing.T) {
		var buf bytes.Buffer
		const call = `InfoContext(ctx, "lines queued", "lines", 1)`
		logger(&buf).InfoContext(withPass(context.Background(), cirrine), "lines queued", "lines", 1)
		rec, _ := decode(t, call, &buf)
		if rec["pass_id"] != cirrine {
			t.Errorf(`%s: "pass_id" = %s, want %q`, call, field(rec, "pass_id"), cirrine)
		}
		if _, ok := rec["request_id"]; ok {
			t.Errorf(`%s with no request ID wrote "request_id" = %s, want none`, call, field(rec, "request_id"))
		}
	})
	t.Run("a context without them", func(t *testing.T) {
		var buf bytes.Buffer
		const call = `InfoContext(context.Background(), "pass planned")`
		logger(&buf).InfoContext(context.Background(), "pass planned")
		rec, raw := decode(t, call, &buf)
		if _, ok := rec["pass_id"]; ok {
			t.Errorf("%s wrote %s, want no pass_id", call, raw)
		}
		if _, ok := rec["request_id"]; ok {
			t.Errorf("%s wrote %s, want no request_id", call, raw)
		}
	})
	t.Run("Info has no context to read", func(t *testing.T) {
		var buf bytes.Buffer
		logger(&buf).Info("pass planned")
		if _, raw := decode(t, `Info("pass planned")`, &buf); strings.Contains(raw, "_id") {
			t.Errorf(`Info("pass planned") wrote %s, want no IDs`, raw)
		}
	})
	t.Run("a DEBUG record is still dropped at INFO", func(t *testing.T) {
		var buf bytes.Buffer
		logger(&buf).DebugContext(booking(), "frame received", "bytes", 128)
		if buf.Len() != 0 {
			t.Errorf(`DebugContext(ctx, "frame received", ...) at INFO wrote %q, want nothing`, buf.String())
		}
	})
}
