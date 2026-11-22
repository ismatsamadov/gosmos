package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// brine3 is the first pass Lantern Atoll flew after the restart at 03:14 on 8 November.
var brine3 = PassView{
	ID: "brine-3/LAN/14384",
	Pass: Pass{
		SatID:      "brine-3",
		Station:    "LAN",
		AOS:        time.Date(2028, 11, 8, 3, 26, 17, 0, time.UTC),
		LOS:        time.Date(2028, 11, 8, 3, 32, 10, 0, time.UTC),
		MaxElevDeg: 16.6,
	},
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

func TestLogFlown(t *testing.T) {
	const call = "logFlown(l, brine-3/LAN/14384, 3)"
	rec := record(t, call, func(l *slog.Logger) { logFlown(l, brine3, 3) })
	if rec["level"] != "INFO" || rec["msg"] != "pass flown" {
		t.Errorf(`%s: "level" = %s and "msg" = %s, want "INFO" and "pass flown"`, call, field(rec, "level"), field(rec, "msg"))
	}
	if rec["pass_id"] != brine3.ID {
		t.Errorf(`%s: "pass_id" = %s, want %q`, call, field(rec, "pass_id"), brine3.ID)
	}
	if rec["max_elev_deg"] != 16.6 {
		t.Errorf(`%s: "max_elev_deg" = %s, want 16.6, a number`, call, field(rec, "max_elev_deg"))
	}
	if rec["lines"] != 3.0 {
		t.Errorf(`%s: "lines" = %s, want 3, a number`, call, field(rec, "lines"))
	}
}
