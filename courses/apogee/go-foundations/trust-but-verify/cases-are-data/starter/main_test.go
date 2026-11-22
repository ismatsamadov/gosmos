package main

import (
	"reflect"
	"testing"
)

func TestEnoughCases(t *testing.T) {
	if len(canonicalCases) < 3 {
		t.Errorf("len(canonicalCases) = %d, want at least 3", len(canonicalCases))
	}
	seen := map[string]bool{}
	for i, c := range canonicalCases {
		switch {
		case c.name == "":
			t.Errorf("case %d has no name; give every case one", i+1)
		case seen[c.name]:
			t.Errorf("two cases are named %q; give each its own name", c.name)
		}
		seen[c.name] = true
		for _, earlier := range canonicalCases[:i] {
			if reflect.DeepEqual(c.cmd, earlier.cmd) {
				t.Errorf("cases %q and %q hold the same command; give each case a different one", earlier.name, c.name)
			}
		}
	}
}

func TestCasesHold(t *testing.T) {
	for _, c := range canonicalCases {
		t.Run(c.name, func(t *testing.T) {
			if got := canonical(c.cmd); got != c.want {
				t.Errorf("canonical(%+v) = %q, want %q", c.cmd, got, c.want)
			}
		})
	}
}
