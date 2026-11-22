package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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

func TestAccessLog(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime})
	stderr := stderrOf(t, func() {
		l := accessLog(h)
		// two requests, written as bookingd's logRequests writes them
		l.Printf("%s %s %s %d %dB %v", "req-000007", "POST", "/bookings", 201, 236, 4*time.Millisecond)
		l.Printf("%s %s %s %d %dB %v", "req-000008", "GET", "/health", 200, 16, time.Duration(0))
	})
	if stderr != "" {
		t.Errorf("accessLog(h) wrote %q to standard error, want nothing there", stderr)
	}
	var recs []map[string]any
	for _, line := range strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n") {
		var rec map[string]any
		if line != "" && json.Unmarshal([]byte(line), &rec) == nil {
			recs = append(recs, rec)
		}
	}
	want := []string{"req-000007 POST /bookings 201 236B 4ms", "req-000008 GET /health 200 16B 0s"}
	if len(recs) != len(want) {
		t.Fatalf("two Printf calls on accessLog(h) sent %d records to h, want 2", len(recs))
	}
	for i, rec := range recs {
		if rec["level"] != "INFO" || rec["msg"] != want[i] {
			t.Errorf(`record %d: "level" = %s and "msg" = %s, want "INFO" and %q`, i+1, field(rec, "level"), field(rec, "msg"), want[i])
		}
	}
}
