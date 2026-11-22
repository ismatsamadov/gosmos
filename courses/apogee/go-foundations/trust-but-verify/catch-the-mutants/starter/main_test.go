package main

import "testing"

// mutants are Mirela's broken copies of accept, each with one small bug.
var mutants = []struct {
	code, bug string
	accept    func(seq, last uint64) bool
}{
	{"seq >= last", "accepts a replay", func(seq, last uint64) bool { return seq >= last }},
	{"seq > last+1", "rejects the next command", func(seq, last uint64) bool { return seq > last+1 }},
	{"seq != last", "accepts an old command", func(seq, last uint64) bool { return seq != last }},
}

func TestCasesHold(t *testing.T) {
	for _, c := range acceptCases {
		t.Run(c.name, func(t *testing.T) {
			if got := accept(c.seq, c.last); got != c.want {
				t.Errorf("accept(%d, %d) = %v, want %v", c.seq, c.last, got, c.want)
			}
		})
	}
}

func TestCasesCatchMutants(t *testing.T) {
	for _, m := range mutants {
		caught := false
		for _, c := range acceptCases {
			if m.accept(c.seq, c.last) != c.want {
				caught = true
			}
		}
		if !caught {
			t.Errorf("mutant %q (%s) survived: it gets every case right, want at least one wrong", m.code, m.bug)
		}
	}
}
