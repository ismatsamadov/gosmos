package main

import (
	"fmt"
	"testing"
	"time"
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

// jossLines are fourteen ways Joss's next flight software might write Saturday's frame 2214,
// with the battery reading 7.91 V: lines Thistle-2 should never send, and tlmd must refuse.
var jossLines = []string{
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=+7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=07.91 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12+00:00 batt_v=7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 02214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=7.9 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1e0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=NaN temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=safe",
	"Thistle-2 2214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 -1 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12.000Z batt_v=7.91 temp_c=1.0 mode=SAFE",
	"thistle-2  2214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=SAFE",
	"thistle-2 2214 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=SAFE crc=6e",
}

// FuzzParseLine runs checkLine on parseLine for every seed in lineSeeds, each of Joss's lines and
// every file under testdata/fuzz/FuzzParseLine: Saturday's lines and the fuzzer's finds.
func FuzzParseLine(f *testing.F) {
	for _, line := range append(lineSeeds, jossLines...) {
		f.Add(line)
	}
	f.Fuzz(func(t *testing.T, line string) {
		checkLine(t, safely(t, parseLine), line)
	})
}

func TestJossLines(t *testing.T) {
	var accepted []string
	for _, line := range jossLines {
		if _, err := safely(t, parseLine)(line); err == nil {
			accepted = append(accepted, line)
		}
	}
	if len(accepted) > 0 {
		t.Errorf("parseLine accepted %d of Joss's %d lines, the first %q; want every one refused", len(accepted), len(jossLines), accepted[0])
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name, line, want string
	}{
		{"a sequence number of -1", "thistle-2 -1 2029-04-14T03:36:12Z batt_v=7.91 temp_c=1.0 mode=SAFE",
			`line: sequence "-1": want a whole number from 0 to 4294967295`},
		{"a time ending +00:00", "thistle-2 2214 2029-04-14T03:36:12+00:00 batt_v=7.91 temp_c=1.0 mode=SAFE",
			`line: time "2029-04-14T03:36:12+00:00": want 2006-01-02T15:04:05Z`},
		{"a battery of 07.91", "thistle-2 2214 2029-04-14T03:36:12Z batt_v=07.91 temp_c=1.0 mode=SAFE",
			`line: not canonical`},
		{"the mode in lowercase", "cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=nominal",
			`line: mode "nominal": want NOMINAL, SAFE, DETUMBLE or LEOP`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := safely(t, parseLine)(tt.line); err == nil || err.Error() != tt.want {
				t.Errorf("parseLine(%q) error = %v, want %q", tt.line, err, tt.want)
			}
		})
	}
}

// frames returns n of sat's telemetry lines, one every 30 seconds from at, numbered from seq.
func frames(sat string, seq int, at time.Time, n int) []string {
	lines := make([]string, n)
	for i := range lines {
		t := at.Add(time.Duration(i) * 30 * time.Second)
		lines[i] = fmt.Sprintf("%s %d %s batt_v=8.01 temp_c=10.5 mode=NOMINAL", sat, seq+i, t.Format(time.RFC3339))
	}
	return lines
}

func TestDrillLines(t *testing.T) {
	// Ten lines from each of five passes of 7 March, the drill's day, numbered as the drill's feeds number them.
	var lines []string
	for _, p := range []struct {
		sat, at string
		seq     int
	}{
		{"pipit-2", "08:15:35", 12131}, {"pipit-3", "08:33:05", 7866}, {"furrow-a", "08:32:05", 31464},
		{"furrow-b", "09:19:35", 31539}, {"brine-3", "09:23:05", 24266},
	} {
		at, _ := time.Parse(time.DateTime, "2029-03-07 "+p.at)
		lines = append(lines, frames(p.sat, p.seq, at, 10)...)
	}
	for _, line := range lines {
		f, err := safely(t, parseLine)(line)
		if err != nil {
			t.Fatalf("parseLine(%q) error = %v, want nil", line, err)
		}
		if got := formatLine(f); got != line {
			t.Fatalf("formatLine(parseLine(%q)) = %q, want the line back", line, got)
		}
	}
}
