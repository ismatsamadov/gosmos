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

// noSatellite is the reference parser that splits a line at single spaces, so a line starting
// with one has an empty satellite ID.
func noSatellite(line string) (TelemetryFrame, error) {
	rest, ok := strings.CutPrefix(line, " ")
	if !ok {
		return reference(line)
	}
	f, err := reference("x " + rest)
	f.SatID = ""
	return f, err
}

// mutants are Mirela's broken copies of her reference parser, each with one bug planted.
var mutants = []struct {
	bug   string
	parse func(string) (TelemetryFrame, error)
}{
	{"accepts batt_v=nan", planted(" batt_v=nan ", " batt_v=0.00 ", func(f *TelemetryFrame) { f.BattV = math.NaN() })},
	{"accepts temp_c=-Inf", planted(" temp_c=-Inf ", " temp_c=0.0 ", func(f *TelemetryFrame) { f.TempC = math.Inf(-1) })},
	{"accepts mode=nominal", planted(" mode=nominal", " mode=NOMINAL", func(f *TelemetryFrame) { f.Mode = "nominal" })},
	{"accepts a line with no satellite ID", noSatellite},
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

// FuzzParseLine runs checkLine on every seed in lineSeeds and every file under
// testdata/fuzz/FuzzParseLine, Saturday's lines from Thistle-2.
func FuzzParseLine(f *testing.F) {
	for _, line := range lineSeeds {
		f.Add(line)
	}
	parse := parser(f)
	f.Fuzz(func(t *testing.T, line string) {
		checkLine(t, safely(t, parse), line)
	})
}

func TestChecksCatchMutants(t *testing.T) {
	for _, m := range mutants {
		t.Run("the mutant that "+m.bug, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^FuzzParseLine$", "-test.timeout=5s")
			cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
			err := cmd.Run()
			var exit *exec.ExitError
			switch {
			case err == nil:
				t.Errorf("the mutant that %s survived: checkLine reported nothing on any seed", m.bug)
			case !errors.As(err, &exit):
				t.Fatalf("running the mutant that %s: %v", m.bug, err)
			}
		})
	}
}

func TestRefusedLines(t *testing.T) {
	// refuseAll refuses every line, and hands back a frame no satellite could send beside its error.
	refuseAll := func(string) (TelemetryFrame, error) {
		return TelemetryFrame{BattV: math.NaN(), Mode: "?"}, errors.New("line: refused, for a reason of its own")
	}
	passed := t.Run("a parser that refuses every seed", func(t *testing.T) {
		for _, line := range lineSeeds {
			checkLine(t, refuseAll, line)
		}
	})
	if !passed {
		t.Errorf("checkLine failed lines the parser refused, want it to check a frame only when the error is nil")
	}
}
