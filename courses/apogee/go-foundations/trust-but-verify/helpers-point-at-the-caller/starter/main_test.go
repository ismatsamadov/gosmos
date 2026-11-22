package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"testing"
)

func TestAssertNear(t *testing.T) {
	assertNear(t, dopplerHz(437.25e6, -7000), 10209.56, 0.01)
	assertNear(t, dopplerHz(437.25e6, 7000), -10209.56, 0.01)
	assertNear(t, dopplerHz(2245e6, -6500), 48675.34, 0.01)
}

// TestProbe fails on purpose, but only when APOGEE_PROBE is set: TestHelperPointsAtCaller
// runs it in a fresh copy of this test binary to see where the failure points.
func TestProbe(t *testing.T) {
	if os.Getenv("APOGEE_PROBE") == "" {
		return
	}
	_, _, line, _ := runtime.Caller(0)
	fmt.Printf("probe line %d\n", line+2) // the assertNear call, two lines down
	assertNear(t, dopplerHz(437.25e6, -7000), 10200, 1)
}

func TestHelperPointsAtCaller(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestProbe$", "-test.timeout=5s")
	cmd.Env = append(os.Environ(), "APOGEE_PROBE=1")
	out, _ := cmd.CombinedOutput()
	probe := regexp.MustCompile(`probe line (\d+)`).FindSubmatch(out)
	if probe == nil {
		t.Fatalf("the probe didn't run")
	}
	line, _ := strconv.Atoi(string(probe[1]))
	failure := regexp.MustCompile(`(?m)^\s+(\w+\.go:\d+): `).FindSubmatch(out)
	if failure == nil {
		t.Fatalf("assertNear(t, 10209.56, 10200, 1) reported no failure; it must fail outside the tolerance")
	}
	if got, want := string(failure[1]), fmt.Sprintf("main_test.go:%d", line); got != want {
		t.Errorf("a failing assertNear points at %s; want %s, the line that called it", got, want)
	}
}
