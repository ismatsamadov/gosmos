package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// lines splits what a logger wrote into its lines.
func lines(out string) []string {
	if out == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(out, "\n"), "\n")
}

// object decodes one line as a JSON object, and stops the test, quoting the line, when it isn't one.
func object(t *testing.T, line string) map[string]any {
	t.Helper()
	var rec map[string]any
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		t.Fatalf("newLogger(w) wrote %q, want a JSON object (%v)", line, err)
	}
	return rec
}

// show writes a decoded value as JSON writes it, or says it is missing.
func show(v any) string {
	if v == nil {
		return "missing"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func TestOneRecordOneObject(t *testing.T) {
	var buf bytes.Buffer
	newLogger(&buf).Error("pass failed", "site", "LAN")
	got := lines(buf.String())
	if len(got) != 1 {
		t.Fatalf(`newLogger(w).Error("pass failed", "site", "LAN") wrote %d lines, want 1`, len(got))
	}
	rec := object(t, got[0])
	for _, f := range []struct{ key, want string }{
		{"level", "ERROR"},
		{"msg", "pass failed"},
		{"site", "LAN"},
	} {
		if rec[f.key] != f.want {
			t.Errorf("%q = %s, want %q", f.key, show(rec[f.key]), f.want)
		}
	}
}

func TestTwoRecordsTwoLines(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf)
	l.Warn("carrier lost", "site", "LAN", "margin_db", -1.5)
	l.Error("pass failed", "site", "LAN")
	got := lines(buf.String())
	if len(got) != 2 {
		t.Fatalf("two records wrote %d lines, want 2: %q", len(got), buf.String())
	}
	for i, want := range []string{"carrier lost", "pass failed"} {
		if rec := object(t, got[i]); rec["msg"] != want {
			t.Errorf(`line %d: "msg" = %s, want %q`, i+1, show(rec["msg"]), want)
		}
	}
}
