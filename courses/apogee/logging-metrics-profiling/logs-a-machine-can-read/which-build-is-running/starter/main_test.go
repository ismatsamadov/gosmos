package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
)

// field writes the value of rec's key as JSON writes it, or says it is missing.
func field(rec map[string]any, key string) string {
	v, ok := rec[key]
	if !ok {
		return "missing"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// TestLogStart reads the start line back and compares it with the build info of this test
// binary, which go test builds from the lesson's module: its Go release, its module's path,
// and the version (devel), since go test never takes one from version control.
func TestLogStart(t *testing.T) {
	var buf bytes.Buffer
	logStart(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime})))
	const call = "logStart(l)"
	out := strings.TrimSuffix(buf.String(), "\n")
	if n := strings.Count(buf.String(), "\n"); n != 1 {
		t.Fatalf("%s wrote %d lines, want 1", call, n)
	}
	var rec map[string]any
	if err := json.Unmarshal([]byte(out), &rec); err != nil {
		t.Fatalf("%s wrote %q, want a JSON object (%v)", call, out, err)
	}
	if rec["level"] != "INFO" || rec["msg"] != "program started" {
		t.Errorf(`%s: "level" = %s and "msg" = %s, want "INFO" and "program started"`, call, field(rec, "level"), field(rec, "msg"))
	}
	if want := runtime.Version(); rec["go_version"] != want {
		t.Errorf(`%s: "go_version" = %s, want %q, the Go release that built this test`, call, field(rec, "go_version"), want)
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		t.Fatal("debug.ReadBuildInfo() in the test = no build info, want this test binary's")
	}
	switch rec["module"] {
	case bi.Main.Path:
	case bi.Path:
		t.Errorf(`%s: "module" = %q, want %q: %q is the main package's path, bi.Path, and under go test it ends in .test`, call, bi.Path, bi.Main.Path, bi.Path)
	default:
		t.Errorf(`%s: "module" = %s, want %q, the main module's path`, call, field(rec, "module"), bi.Main.Path)
	}
	if rec["module_version"] != bi.Main.Version {
		t.Errorf(`%s: "module_version" = %s, want %q, the main module's version`, call, field(rec, "module_version"), bi.Main.Version)
	}
}
