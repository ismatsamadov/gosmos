package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// stowsAbove72 is windCall with > where >= belongs, so it tracks at exactly 72 km/h.
func stowsAbove72(kmh float64) string {
	switch {
	case kmh > 72:
		return "stow"
	case kmh >= 50:
		return "track with warning"
	default:
		return "track"
	}
}

// warnsAbove50 is windCall with > where >= belongs, so it gives no warning at 50 km/h.
func warnsAbove50(kmh float64) string {
	switch {
	case kmh >= 72:
		return "stow"
	case kmh > 50:
		return "track with warning"
	default:
		return "track"
	}
}

// mutants are broken copies of windCall, each with one planted bug.
var mutants = []struct {
	bug      string
	windCall func(kmh float64) string
}{
	{"keeps tracking at 72 km/h", stowsAbove72},
	{"gives no warning at 50 km/h", warnsAbove50},
}

func TestWind(t *testing.T) {
	checkWind(t, windCall)
}

func TestSubtestsNamed(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestWind$", "-test.v", "-test.timeout=5s")
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		t.Fatalf("running TestWind: %v", err)
	}
	for _, c := range windCases {
		name := "TestWind/" + strings.ReplaceAll(c.name, " ", "_")
		if !strings.Contains(string(out), "=== RUN   "+name+"\n") {
			t.Errorf("TestWind ran no subtest %s: run each case with t.Run(c.name, ...)", name)
		}
	}
}

// TestMutant runs checkWind against the mutant APOGEE_MUTANT names. In a normal run it does
// nothing: TestChecksCatchMutants runs it in a fresh copy of this test binary.
func TestMutant(t *testing.T) {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return
	}
	for _, m := range mutants {
		if m.bug == bug {
			checkWind(t, m.windCall)
			return
		}
	}
	t.Fatalf("no mutant %q", bug)
}

// TestChecksCatchMutants runs each mutant in a fresh copy of this test binary. The run must
// fail, and the failure must belong to a subtest, so that the log names the case.
func TestChecksCatchMutants(t *testing.T) {
	for _, m := range mutants {
		cmd := exec.Command(os.Args[0], "-test.run=^TestMutant$", "-test.timeout=5s")
		cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
		out, err := cmd.CombinedOutput()
		var exit *exec.ExitError
		switch {
		case err == nil:
			t.Errorf("the mutant that %s survived: checkWind reported nothing against it", m.bug)
		case !errors.As(err, &exit):
			t.Fatalf("running the mutant that %s: %v", m.bug, err)
		case !strings.Contains(string(out), "--- FAIL: TestMutant/"):
			t.Errorf("the mutant that %s failed no subtest, only TestMutant; want the failing case's subtest to fail: report with the t that t.Run hands you", m.bug)
		}
	}
}
