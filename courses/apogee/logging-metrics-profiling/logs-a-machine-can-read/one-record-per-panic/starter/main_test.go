package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// orbitOf returns the orbit number at the end of a pass ID; an ID without one panics, as on the bench.
func orbitOf(id string) string {
	return strings.Split(id, "/")[2]
}

// markQueued records pass id in queued, and panics when nobody made the map.
func markQueued(queued map[string]bool, id string) {
	queued[id] = true
}

// panicText runs f, which must panic, and returns the panic's value as text.
func panicText(f func()) (text string) {
	defer func() {
		if v := recover(); v != nil {
			text = fmt.Sprint(v)
		}
	}()
	f()
	return "no panic"
}

// stderrOf runs f with standard error sent to a file of the test's own, and returns what f wrote there.
func stderrOf(t *testing.T, f func()) string {
	t.Helper()
	file, err := os.Create(filepath.Join(t.TempDir(), "stderr.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	saved := os.Stderr
	os.Stderr = file
	defer func() { os.Stderr = saved }()
	f()
	b, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
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

// short cuts s to its first n bytes, to quote a long value in a message.
func short(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

func TestPanicIsOneRecord(t *testing.T) {
	tests := []struct {
		name string
		id   string
		work func()
		fn   string // the function that panics
	}{
		{"the bench's pass ID", "cirrine-1/LAN", func() { _ = orbitOf("cirrine-1/LAN") }, "orbitOf"},
		{"a map nobody made", "brine-3/LAN/14384", func() { markQueued(nil, "brine-3/LAN/14384") }, "markQueued"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := fmt.Sprintf("pass %s: panic: %s", tt.id, panicText(tt.work))
			call := fmt.Sprintf("guardPass(l, %q, work)", tt.id)
			var buf bytes.Buffer
			l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime}))
			stderr := stderrOf(t, func() { guardPass(l, tt.id, tt.work) })
			if stderr != "" {
				t.Errorf("%s wrote %d bytes to standard error, beginning %q, want nothing there", call, len(stderr), short(stderr, 30))
			}
			out := strings.TrimSuffix(buf.String(), "\n")
			if n := strings.Count(buf.String(), "\n"); n != 1 {
				t.Fatalf("%s wrote %d lines, want 1 record", call, n)
			}
			var rec map[string]any
			if err := json.Unmarshal([]byte(out), &rec); err != nil {
				t.Fatalf("%s wrote %q, want a JSON object (%v)", call, short(out, 60), err)
			}
			if rec["level"] != "ERROR" || rec["msg"] != "pass failed" {
				t.Errorf(`%s: "level" = %s and "msg" = %s, want "ERROR" and "pass failed"`, call, field(rec, "level"), field(rec, "msg"))
			}
			if rec["pass_id"] != tt.id {
				t.Errorf(`%s: "pass_id" = %s, want %q`, call, field(rec, "pass_id"), tt.id)
			}
			if rec["err"] != want {
				t.Errorf(`%s: "err" = %s, want %q`, call, short(field(rec, "err"), 90), want)
			}
			stack, ok := rec["stack"].(string)
			switch {
			case !ok:
				t.Errorf(`%s: "stack" = %s, want the stack as text`, call, short(field(rec, "stack"), 40))
			case !strings.HasPrefix(stack, "goroutine "):
				t.Errorf(`%s: "stack" = %q, want text beginning "goroutine ", as debug.Stack writes it`, call, short(stack, 30))
			case !strings.Contains(stack, "."+tt.fn+"("):
				t.Errorf(`%s: "stack" names no %s, want the stack of the panic, taken in the deferred function`, call, tt.fn)
			}
		})
	}
}

func TestNoPanicNoRecord(t *testing.T) {
	const call = `guardPass(l, "cirrine-1/LAN/16402", work)`
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime}))
	ran := false
	guardPass(l, "cirrine-1/LAN/16402", func() { ran = true })
	if !ran {
		t.Errorf("%s never ran work, want it run once", call)
	}
	if buf.Len() != 0 {
		t.Errorf("%s wrote %q for a pass that didn't panic, want nothing", call, short(buf.String(), 60))
	}
}
