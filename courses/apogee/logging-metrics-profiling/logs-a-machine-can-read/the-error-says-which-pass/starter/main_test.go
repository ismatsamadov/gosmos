package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

// deadline is the error a pass ends with when its time runs out at LOS: its ID, then the context's.
func deadline(id string) error {
	return fmt.Errorf("pass %s: %w", id, context.DeadlineExceeded)
}

// record calls write with a JSON logger, like the ones the night shift's search reads, and
// decodes the line it writes; it stops the test, naming call, unless there is exactly one.
func record(t *testing.T, call string, write func(l *slog.Logger)) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	write(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime})))
	out := strings.TrimSuffix(buf.String(), "\n")
	if out == "" {
		t.Fatalf("%s wrote nothing, want one line", call)
	}
	if n := strings.Count(out, "\n") + 1; n != 1 {
		t.Fatalf("%s wrote %d lines, want 1: %q", call, n, buf.String())
	}
	var rec map[string]any
	if err := json.Unmarshal([]byte(out), &rec); err != nil {
		t.Fatalf("%s wrote %q, want a JSON object (%v)", call, out, err)
	}
	return rec
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

// TestLogResult replays how Lantern Atoll's passes of 8 November ended.
func TestLogResult(t *testing.T) {
	night := []struct {
		name string
		id   string
		err  error
	}{
		{"Brine-3 at 01:54:41", "brine-3/LAN/14383", deadline("brine-3/LAN/14383")},
		{"Cirrine-1 at 02:02:39", "cirrine-1/LAN/16402", deadline("cirrine-1/LAN/16402")},
		{"Furrow-A at 02:46:24", "furrow-a/LAN/20998", deadline("furrow-a/LAN/20998")},
		{"Brine-3 after the restart", "brine-3/LAN/14384", nil},
	}
	for _, tt := range night {
		t.Run(tt.name, func(t *testing.T) {
			call := fmt.Sprintf("logResult(l, %q, err)", tt.id)
			if tt.err == nil {
				call = fmt.Sprintf("logResult(l, %q, nil)", tt.id)
			}
			rec := record(t, call, func(l *slog.Logger) { logResult(l, tt.id, tt.err) })
			if rec["pass_id"] != tt.id {
				t.Errorf(`%s: "pass_id" = %s, want %q`, call, field(rec, "pass_id"), tt.id)
			}
			if tt.err == nil {
				if rec["level"] != "INFO" || rec["msg"] != "pass flown" {
					t.Errorf(`%s: "level" = %s and "msg" = %s, want "INFO" and "pass flown"`, call, field(rec, "level"), field(rec, "msg"))
				}
				if _, ok := rec["err"]; ok {
					t.Errorf(`%s wrote "err":%s, want no "err" key`, call, field(rec, "err"))
				}
				return
			}
			if rec["level"] != "ERROR" || rec["msg"] != "pass failed" {
				t.Errorf(`%s: "level" = %s and "msg" = %s, want "ERROR" and "pass failed"`, call, field(rec, "level"), field(rec, "msg"))
			}
			if want := tt.err.Error(); rec["err"] != want {
				t.Errorf(`%s: "err" = %s, want %q, the whole text`, call, field(rec, "err"), want)
			}
			if !errors.Is(tt.err, context.DeadlineExceeded) {
				t.Errorf("errors.Is(err, context.DeadlineExceeded) = false after %s, want true", call)
			}
		})
	}
}
