package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func BenchmarkWindows(b *testing.B) { benchWindows(b) }

// procs is the -N a benchmark's line adds to its name.
var procs = regexp.MustCompile(`-\d+$`)

// windowsLines runs BenchmarkWindows again in a new run of this test
// binary, at 2 turns each and on one CPU, as AllocsPerRun counts, and
// returns its result lines by name, each split into fields after the
// name.
func windowsLines(t *testing.T) map[string][]string {
	limit := 5 * time.Second // real work: a few tenths of a second
	if raceOn {
		limit = time.Minute
	}
	cmd := exec.Command(os.Args[0], "-test.run=^$",
		"-test.bench=^BenchmarkWindows$", "-test.benchtime=2x",
		"-test.cpu=1", "-test.timeout="+limit.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		first := strings.SplitN(string(out), "\n", 12)
		t.Fatalf("go test -bench=Windows failed: %v\n%s", err, strings.Join(first[:min(len(first), 11)], "\n"))
	}
	lines := map[string][]string{}
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) > 1 && strings.HasPrefix(f[0], "Benchmark") {
			lines[procs.ReplaceAllString(f[0], "")] = f[1:]
		}
	}
	return lines
}

// allocsPerOp returns the number before "allocs/op" in a line's
// fields, and false when the line has none.
func allocsPerOp(fields []string) (int64, bool) {
	for i := 1; i < len(fields); i++ {
		if fields[i] == "allocs/op" {
			n, err := strconv.ParseInt(fields[i-1], 10, 64)
			return n, err == nil
		}
	}
	return 0, false
}

func TestBenchWindows(t *testing.T) {
	type sub struct {
		name  string
		plan  func(*Network, Window, Config) []Contact
		hours int
	}
	var subs []sub
	for _, v := range []string{"old", "new"} {
		plan := planOld
		if v == "new" {
			plan = planNew
		}
		for _, h := range []int{6, 12, 24} {
			name := fmt.Sprintf("BenchmarkWindows/%s/hours=%d", v, h)
			subs = append(subs, sub{name, plan, h})
		}
	}

	lines := windowsLines(t)
	var names []string
	for name := range lines {
		names = append(names, name)
	}
	slices.Sort(names)
	missing := 0
	for _, s := range subs {
		if _, ok := lines[s.name]; !ok {
			missing++
		}
	}
	if missing == len(subs) {
		t.Fatalf("go test -bench=Windows listed %s, want BenchmarkWindows/old and /new, each with hours=6, hours=12 and hours=24", strings.Join(names, ", "))
	}
	for _, name := range names {
		if !slices.ContainsFunc(subs, func(s sub) bool { return s.name == name }) {
			t.Errorf("go test -bench=Windows listed %s, want only old and new, each with hours=6, hours=12 and hours=24", name)
		}
	}
	for _, s := range subs {
		f, ok := lines[s.name]
		switch {
		case !ok:
			t.Errorf("go test -bench=Windows listed no %s, want it", s.name)
		case f[0] != "2":
			t.Errorf("%s ran %s turns, want 2, as -benchtime=2x asks", s.name, f[0])
		default:
			if _, ok := allocsPerOp(f); !ok {
				t.Errorf("%s reported no allocs/op, want them: b.ReportAllocs in each sub-benchmark", s.name)
			}
		}
	}
	if t.Failed() {
		return
	}

	t.Run("each measures its planner on its window", func(t *testing.T) {
		if raceOn {
			t.Skip("allocation counts differ under -race")
		}
		n := newNetwork()
		cfg := Config{MinElevDeg: 10, Turnaround: 2 * time.Minute}
		friday := time.Date(2029, 5, 18, 0, 0, 0, 0, time.UTC)
		planned := map[int][]Contact{}
		for _, s := range subs {
			w := Window{From: friday, To: friday.Add(time.Duration(s.hours) * time.Hour)}
			var contacts []Contact
			own := int64(testing.AllocsPerRun(1, func() {
				contacts = s.plan(n, w, cfg)
			}))
			got, _ := allocsPerOp(lines[s.name])
			if math.Abs(float64(got-own)) > 0.015*float64(own) {
				t.Errorf("%s made %d allocs/op, want %d, as its planner makes for the %d hours from midnight on 18 May", s.name, got, own, s.hours)
			}
			same := func(a, b Contact) bool {
				return a.ID == b.ID && a.Dish == b.Dish
			}
			earlier, ok := planned[s.hours]
			if ok && !slices.EqualFunc(earlier, contacts, same) {
				t.Fatalf("planOld and planNew disagree on the %d hours from midnight on 18 May, want the same contacts", s.hours)
			}
			planned[s.hours] = contacts
		}
	})
}
