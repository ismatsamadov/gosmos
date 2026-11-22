package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

// fullAbove82 is batteryBand with > where >= belongs, so exactly 8.2 V reads OK.
func fullAbove82(v float64) string {
	switch {
	case v > 8.2:
		return "FULL"
	case v >= 7.2:
		return "OK"
	case v >= 6.6:
		return "LOW"
	default:
		return "CRITICAL"
	}
}

// okFrom70 is batteryBand with OK starting at 7.0 V instead of 7.2 V.
func okFrom70(v float64) string {
	switch {
	case v >= 8.2:
		return "FULL"
	case v >= 7.0:
		return "OK"
	case v >= 6.6:
		return "LOW"
	default:
		return "CRITICAL"
	}
}

// lowAbove66 is batteryBand with > where >= belongs, so exactly 6.6 V reads CRITICAL.
func lowAbove66(v float64) string {
	switch {
	case v >= 8.2:
		return "FULL"
	case v >= 7.2:
		return "OK"
	case v > 6.6:
		return "LOW"
	default:
		return "CRITICAL"
	}
}

// mutants are broken copies of batteryBand, each with one planted bug.
var mutants = []struct {
	bug  string
	band func(v float64) string
}{
	{"reads 8.2 V as OK", fullAbove82},
	{"starts OK at 7.0 V", okFrom70},
	{"reads 6.6 V as CRITICAL", lowAbove66},
}

func TestRealCodePasses(t *testing.T) {
	checkBand(t, batteryBand)
}

// TestMutant runs checkBand against the mutant APOGEE_MUTANT names. In a normal run it does
// nothing: TestChecksCatchMutants runs it in a fresh copy of this test binary.
func TestMutant(t *testing.T) {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return
	}
	for _, m := range mutants {
		if m.bug == bug {
			checkBand(t, m.band)
			return
		}
	}
	t.Fatalf("no mutant %q", bug)
}

func TestChecksCatchMutants(t *testing.T) {
	for _, m := range mutants {
		cmd := exec.Command(os.Args[0], "-test.run=^TestMutant$", "-test.timeout=5s")
		cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
		err := cmd.Run()
		var exit *exec.ExitError
		switch {
		case err == nil:
			t.Errorf("the mutant that %s survived: checkBand reported nothing against it", m.bug)
		case !errors.As(err, &exit):
			t.Fatalf("running the mutant that %s: %v", m.bug, err)
		}
	}
}
