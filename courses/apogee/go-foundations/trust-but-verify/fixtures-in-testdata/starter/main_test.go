package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLoadsFixture(t *testing.T) {
	lines := loadFixture(t, "thistle-0511.tlm")
	if len(lines) != 12 {
		t.Fatalf("len(loadFixture(t, %q)) = %d, want 12", "thistle-0511.tlm", len(lines))
	}
	first := "thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL"
	last := "thistle 4373 2027-05-11T05:48:37Z batt_v=7.72 temp_c=2.9 mode=NOMINAL"
	if lines[0] != first || lines[11] != last {
		t.Errorf("first and last lines = %q, %q, want %q, %q", lines[0], lines[11], first, last)
	}
}

func TestBatteryOnFixture(t *testing.T) {
	rejected := 0
	for _, line := range loadFixture(t, "thistle-0511.tlm") {
		if _, err := battery(line); err != nil {
			rejected++
		}
	}
	if rejected != 1 {
		t.Errorf("battery rejected %d fixture lines, want 1: the frame that reads 7.9l", rejected)
	}
}

// TestMissingProbe asks for a fixture that doesn't exist, but only when APOGEE_PROBE is set:
// TestMissingFixtureIsFatal runs it in a fresh copy of this test binary.
func TestMissingProbe(t *testing.T) {
	if os.Getenv("APOGEE_PROBE") == "" {
		return
	}
	loadFixture(t, "thistle-0512.tlm")
	fmt.Println("probe still running")
}

func TestMissingFixtureIsFatal(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestMissingProbe$", "-test.timeout=5s")
	cmd.Env = append(os.Environ(), "APOGEE_PROBE=1")
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		t.Fatalf("running the probe: %v", err)
	}
	if err == nil {
		t.Fatalf("loadFixture(t, %q) let the test pass; a missing fixture must fail it", "thistle-0512.tlm")
	}
	text := string(out)
	switch {
	case strings.Contains(text, "panic:"):
		t.Errorf("loadFixture panicked on a missing file; stop the test with t.Fatalf instead")
	case strings.Contains(text, "probe still running"):
		t.Errorf("the test kept running after loadFixture failed; stop it with t.Fatalf")
	case !strings.Contains(text, "--- FAIL: TestMissingProbe"):
		t.Errorf("the probe exited without failing its test; want loadFixture to fail it with t.Fatalf, not end the program")
	case !strings.Contains(text, "thistle-0512.tlm"):
		t.Errorf("the failure message doesn't name the missing file, thistle-0512.tlm")
	}
}
