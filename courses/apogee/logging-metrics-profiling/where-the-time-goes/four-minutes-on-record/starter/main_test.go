package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"testing"
	"time"
)

// replan is what passplan planned at 05:10 on Tuesday 16 January, once Cape Wren had stowed: the
// whole fleet over Cape Wren and Lantern Atoll, for the rest of the day.
func replan() {
	from := time.Date(2029, 1, 16, 5, 10, 0, 0, time.UTC)
	to := time.Date(2029, 1, 17, 0, 0, 0, 0, time.UTC)
	planDay(fleet(), elements(), []Station{station("CWR"), station("LAN")}, nil, from, to)
}

func TestProfileToMissingFolder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "cpu.prof")
	ran := false
	err := profileTo(path, func() { ran = true })
	if err == nil {
		t.Errorf("profileTo(a path in a missing folder, work) = <nil>, want the error from creating the file")
	}
	if ran {
		t.Errorf("profileTo(a path in a missing folder, work) ran work, want it left unrun")
	}
}

func TestProfileToWhileAnotherRuns(t *testing.T) {
	if err := pprof.StartCPUProfile(io.Discard); err != nil {
		t.Fatalf("pprof.StartCPUProfile(io.Discard) before profileTo = %v, want <nil>", err)
	}
	defer pprof.StopCPUProfile()
	ran := false
	err := profileTo(filepath.Join(t.TempDir(), "cpu.prof"), func() { ran = true })
	if err == nil {
		t.Errorf("profileTo(path, work) with another profile running = <nil>, want StartCPUProfile's error")
	}
	if ran {
		t.Errorf("profileTo(path, work) with another profile running ran work, want it left unrun")
	}
}

func TestProfileTo(t *testing.T) {
	os.Remove("cpu.prof")
	ran := false
	if err := profileTo("cpu.prof", func() { ran = true; replan() }); err != nil {
		t.Fatalf(`profileTo("cpu.prof", work) = %v, want <nil>`, err)
	}
	if !ran {
		t.Errorf(`profileTo("cpu.prof", work) never ran work`)
	}
	b, err := os.ReadFile("cpu.prof")
	switch {
	case err != nil:
		t.Errorf(`after profileTo("cpu.prof", work), reading cpu.prof gives %v, want the profile`, err)
	case len(b) == 0:
		t.Errorf("cpu.prof is empty, want a profile: it is written when the profile stops")
	case len(b) < 2 || b[0] != 0x1f || b[1] != 0x8b:
		t.Errorf("cpu.prof starts with % x, want 1f 8b, the first bytes of a profile", b[:min(len(b), 2)])
	}
	if err := pprof.StartCPUProfile(io.Discard); err != nil {
		t.Errorf("pprof.StartCPUProfile(io.Discard) after profileTo = %v, want <nil>: profileTo's profile is still running", err)
		return
	}
	pprof.StopCPUProfile()
}
