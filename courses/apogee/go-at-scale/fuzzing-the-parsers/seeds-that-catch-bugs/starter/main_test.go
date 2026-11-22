package main

import (
	"errors"
	"math"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// safely returns parse wrapped so that a panic fails the test as parseLine("...") panicked: <value>, instead of ending the test binary.
func safely(t *testing.T, parse func(string) (TelemetryFrame, error)) func(string) (TelemetryFrame, error) {
	return func(line string) (TelemetryFrame, error) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parseLine(%q) panicked: %v", line, r)
			}
		}()
		return parse(line)
	}
}

// reference, Mirela's reference parser, is in a file the lesson keeps hidden. It accepts exactly
// the lines formatLine writes, with readings in range and a known mode, and refuses the rest.

// planted returns the reference parser with one bug planted: it reads a line holding bad as if it
// held good, then spoil changes the frame it accepts.
func planted(bad, good string, spoil func(*TelemetryFrame)) func(string) (TelemetryFrame, error) {
	return func(line string) (TelemetryFrame, error) {
		if !strings.Contains(line, bad) {
			return reference(line)
		}
		f, err := reference(strings.Replace(line, bad, good, 1))
		if err == nil {
			spoil(&f)
		}
		return f, err
	}
}

// fiveFields is the reference parser with an off-by-one count: it lets a line of five fields
// through, then reads its sixth.
func fiveFields(line string) (TelemetryFrame, error) {
	if f := strings.Fields(line); len(f) == 5 {
		return TelemetryFrame{SatID: f[0], Mode: Mode(f[5])}, nil
	}
	return reference(line)
}

// mutants are Mirela's broken copies of her reference parser, each with one bug planted.
var mutants = []struct {
	bug   string
	parse func(string) (TelemetryFrame, error)
}{
	{"panics on a line of five fields", fiveFields},
	{"accepts batt_v=nan", planted(" batt_v=nan ", " batt_v=0.00 ", func(f *TelemetryFrame) { f.BattV = math.NaN() })},
	{"accepts mode=nominal", planted(" mode=nominal", " mode=NOMINAL", func(f *TelemetryFrame) { f.Mode = "nominal" })},
}

// parser returns the parser this run checks: the reference parser, or the mutant APOGEE_MUTANT names.
func parser(tb testing.TB) func(string) (TelemetryFrame, error) {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return reference
	}
	for _, m := range mutants {
		if m.bug == bug {
			return m.parse
		}
	}
	tb.Fatalf("no mutant %q", bug)
	return nil
}

// FuzzParseLine checks Mirela's two rules on every seed: the parser never panics, and a line it
// accepts has a battery from 0 to 20 volts and a known mode.
func FuzzParseLine(f *testing.F) {
	for _, line := range lineSeeds {
		f.Add(line)
	}
	parse := parser(f)
	f.Fuzz(func(t *testing.T, line string) {
		fr, err := safely(t, parse)(line)
		if err != nil {
			return
		}
		if !(fr.BattV >= 0 && fr.BattV <= 20) {
			t.Errorf("parseLine(%q).BattV = %.2f, want 0.00 to 20.00", line, fr.BattV)
		}
		switch fr.Mode {
		case "NOMINAL", "SAFE", "DETUMBLE", "LEOP":
		default:
			t.Errorf("parseLine(%q).Mode = %q, want NOMINAL, SAFE, DETUMBLE or LEOP", line, fr.Mode)
		}
	})
}

func TestSeeds(t *testing.T) {
	if len(lineSeeds) < 3 {
		t.Errorf("len(lineSeeds) = %d, want at least 3: one for each mutant", len(lineSeeds))
	}
	seen := map[string]bool{}
	for _, s := range lineSeeds {
		if seen[s] {
			t.Errorf("lineSeeds holds %q twice, want every seed different", s)
		}
		seen[s] = true
	}
}

func TestSeedsCatchMutants(t *testing.T) {
	for _, m := range mutants {
		t.Run("the mutant that "+m.bug, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^FuzzParseLine$", "-test.timeout=5s")
			cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
			err := cmd.Run()
			var exit *exec.ExitError
			switch {
			case err == nil:
				t.Errorf("the mutant that %s survived: no seed in lineSeeds made FuzzParseLine fail", m.bug)
			case !errors.As(err, &exit):
				t.Fatalf("running the mutant that %s: %v", m.bug, err)
			}
		})
	}
}
