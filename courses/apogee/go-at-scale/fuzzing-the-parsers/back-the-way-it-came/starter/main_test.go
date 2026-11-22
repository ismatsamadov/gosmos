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

// jossLines are four ways Joss's next flight software wrote Saturday's frame 2214, with the
// battery reading 7.91 V, as he sent them to Mirela.
var jossLines = []string{
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=+7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=07.91 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12+00:00 batt_v=7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 02214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=SAFE",
}

// reference, Mirela's reference parser, is in a file the lesson keeps hidden. It accepts exactly
// the lines formatLine writes, with readings in range and a known mode, and refuses the rest.

// lenient returns the reference parser with one way of writing a frame let through: it reads a
// line holding bad as if it held good.
func lenient(bad, good string) func(string) (TelemetryFrame, error) {
	return func(line string) (TelemetryFrame, error) {
		return reference(strings.Replace(line, bad, good, 1))
	}
}

// mutants are Mirela's broken copies of her reference parser, each letting one of Joss's lines through.
var mutants = []struct {
	bug   string
	parse func(string) (TelemetryFrame, error)
}{
	{"accepts batt_v=+7.91", lenient(" batt_v=+7.91 ", " batt_v=7.91 ")},
	{"accepts batt_v=07.91", lenient(" batt_v=07.91 ", " batt_v=7.91 ")},
	{"accepts a time ending +00:00", lenient("+00:00 ", "Z ")},
	{"accepts a sequence number with a leading zero", lenient(" 02214 ", " 2214 ")},
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

// FuzzParseLine runs checkLine on every seed in lineSeeds, each of Joss's lines and every file
// under testdata/fuzz/FuzzParseLine.
func FuzzParseLine(f *testing.F) {
	for _, line := range append(lineSeeds, jossLines...) {
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
		for _, line := range append(lineSeeds, jossLines...) {
			checkLine(t, refuseAll, line)
		}
	})
	if !passed {
		t.Errorf("checkLine failed lines the parser refused, want it to check a frame only when the error is nil")
	}
}
