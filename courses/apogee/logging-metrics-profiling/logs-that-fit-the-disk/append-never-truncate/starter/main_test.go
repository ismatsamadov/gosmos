package main

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// Two of Lantern Atoll's lines, one written before a restart and one after it.
const (
	before = `{"time":"2028-12-04T07:59:58Z","level":"INFO","msg":"lines queued","lines":1}` + "\n"
	after  = `{"time":"2028-12-04T08:00:04Z","level":"INFO","msg":"lines queued","lines":2}` + "\n"
)

// writeOnce opens path with openLog, writes line and closes the file, stopping the test on any error.
func writeOnce(t *testing.T, path, line string) {
	t.Helper()
	f, err := openLog(path)
	if err != nil {
		t.Fatalf("openLog(%q) returned %v, want nil", filepath.Base(path), err)
	}
	if _, err := f.WriteString(line); err != nil {
		t.Fatalf("writing to the file openLog returned: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenLogAppends(t *testing.T) {
	path := filepath.Join(t.TempDir(), "uplinkd.log")
	writeOnce(t, path, before)
	writeOnce(t, path, after) // the Monday restart at 08:00
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != before+after {
		t.Errorf("after two openLog calls, a line written after each, uplinkd.log = %q, want %q", got, before+after)
	}
}

func TestOpenLogMode(t *testing.T) {
	defer syscall.Umask(syscall.Umask(0o022)) // the usual umask, whatever the Mac's is
	path := filepath.Join(t.TempDir(), "uplinkd.log")
	writeOnce(t, path, before)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0640 {
		t.Errorf("a new file from openLog has mode %04o, want 0640", got)
	}
}

func TestOpenLogMissingFolder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "uplinkd.log")
	f, err := openLog(path)
	if err == nil {
		f.Close()
		t.Errorf("openLog in a folder that doesn't exist returned no error, want one")
	}
}
