package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"
)

// cirrine is the pass ID of Cirrine-1's pass of 8 November.
const cirrine = "cirrine-1/LAN/16402"

// logger returns a logger on passHandler around a JSON handler that writes to buf, without the time.
func logger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(passHandler{slog.NewJSONHandler(buf, &slog.HandlerOptions{ReplaceAttr: dropTime})})
}

// line returns the one line in buf, decoded and as written; it stops the test, naming call, unless there is exactly one.
func line(t *testing.T, call string, buf *bytes.Buffer) (map[string]any, string) {
	t.Helper()
	raw := strings.TrimSuffix(buf.String(), "\n")
	var rec map[string]any
	if raw == "" || strings.Contains(raw, "\n") || json.Unmarshal([]byte(raw), &rec) != nil {
		t.Fatalf("%s wrote %q, want one JSON line", call, buf.String())
	}
	return rec, raw
}

func TestWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	ctx := withPass(context.Background(), cirrine)
	const call = `l.With("program", "uplinkd").InfoContext(ctx, "pass started")`
	logger(&buf).With("program", "uplinkd").InfoContext(ctx, "pass started")
	rec, raw := line(t, call, &buf)
	if n := strings.Count(raw, `"pass_id":`); n != 1 {
		t.Errorf("%s wrote %s: %d pass_id keys, want exactly 1", call, raw, n)
	}
	if rec["program"] != "uplinkd" {
		t.Errorf("%s wrote %s, want the inner handler to keep \"program\":\"uplinkd\"", call, raw)
	}
}

func TestWithGroup(t *testing.T) {
	var buf bytes.Buffer
	ctx := withPass(context.Background(), cirrine)
	const call = `l.WithGroup("radio").InfoContext(ctx, "carrier lost", "margin_db", -1.5)`
	logger(&buf).WithGroup("radio").InfoContext(ctx, "carrier lost", "margin_db", -1.5)
	rec, raw := line(t, call, &buf)
	if n := strings.Count(raw, `"pass_id":`); n != 1 {
		t.Errorf("%s wrote %s: %d pass_id keys, want exactly 1", call, raw, n)
	}
	if _, ok := rec["radio"].(map[string]any); !ok {
		t.Errorf("%s wrote %s, want the inner handler to open the group \"radio\"", call, raw)
	}
}

// TestSlogtest checks passHandler against the rules every slog handler keeps.
func TestSlogtest(t *testing.T) {
	var buf bytes.Buffer
	h := passHandler{slog.NewJSONHandler(&buf, nil)}
	results := func() []map[string]any {
		var recs []map[string]any
		for _, raw := range strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n") {
			var rec map[string]any
			if err := json.Unmarshal([]byte(raw), &rec); err != nil {
				t.Fatalf("passHandler wrote %q, want JSON lines", raw)
			}
			recs = append(recs, rec)
		}
		return recs
	}
	if err := slogtest.TestHandler(h, results); err != nil {
		t.Errorf("slogtest.TestHandler(passHandler{...}) = %v, want nil", err)
	}
}
