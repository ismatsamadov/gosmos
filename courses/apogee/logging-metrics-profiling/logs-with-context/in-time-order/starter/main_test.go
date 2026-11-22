package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"os"
	"slices"
	"testing"
	"time"
)

// readLog decodes every line of testdata/<name>, one program's log from the bench replay of 8 November.
func readLog(t *testing.T, name string) []Line {
	t.Helper()
	f, err := os.Open("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var log []Line
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l := Line{Raw: sc.Text()}
		if err := json.Unmarshal(sc.Bytes(), &l); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		log = append(log, l)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return log
}

// inOrder returns the lines of logs by time; lines with equal times come in the order of their logs,
// then of their places in their log.
func inOrder(logs [][]Line) []Line {
	type placed struct {
		line     Line
		log, pos int
	}
	var all []placed
	for i, log := range logs {
		for j, l := range log {
			all = append(all, placed{l, i, j})
		}
	}
	slices.SortFunc(all, func(a, b placed) int {
		return cmp.Or(a.line.Time.Compare(b.line.Time), cmp.Compare(a.log, b.log), cmp.Compare(a.pos, b.pos))
	})
	var lines []Line
	for _, p := range all {
		lines = append(lines, p.line)
	}
	return lines
}

// describe writes a line as the messages show it: time, program, site, message and pass.
func describe(l Line) string {
	return l.Time.UTC().Format(time.RFC3339) + " " + l.Program + " " + l.Site + " " + l.Msg + " " + l.PassID
}

func TestMerge(t *testing.T) {
	names := []string{"bookingd.log", "passplan.log", "uplinkd-hq.log", "uplinkd-lan.log"}
	var logs, saved [][]Line
	for _, name := range names {
		log := readLog(t, name)
		logs = append(logs, log)
		saved = append(saved, slices.Clone(log))
	}
	const call = "merge(bookingd, passplan, uplinkd HQ, uplinkd LAN)"
	got := merge(logs...)
	want := inOrder(saved)
	if len(got) != len(want) {
		t.Fatalf("%s returned %d lines, want all %d", call, len(got), len(want))
	}
	for i := range want {
		if got[i].Raw != want[i].Raw {
			t.Errorf("%s line %d = %q, want %q", call, i+1, describe(got[i]), describe(want[i]))
			break
		}
	}
	for i, log := range logs {
		if !slices.Equal(log, saved[i]) {
			t.Errorf("after %s, the log of %s changed, want every input as it was", call, names[i])
		}
	}
}

func TestMergeNothing(t *testing.T) {
	if got := merge(); len(got) != 0 {
		t.Errorf("merge() returned %d lines, want none", len(got))
	}
}
