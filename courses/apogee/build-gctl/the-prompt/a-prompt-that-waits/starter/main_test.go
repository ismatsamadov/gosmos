package main

import (
	"strings"
	"testing"
)

// session runs s on input and returns what it wrote and what Run returned.
func session(s *Shell, input string) (string, error) {
	var out strings.Builder
	err := s.Run(strings.NewReader(input), &out)
	return out.String(), err
}

// check runs input on a new shell and compares what it wrote with want.
func check(t *testing.T, input, want string) {
	t.Helper()
	got, err := session(newShell(), input)
	if got != want {
		t.Errorf("Run(%q) wrote %q, want %q", input, got, want)
	}
	if err != nil {
		t.Errorf("Run(%q) = %v, want nil", input, err)
	}
}

func TestPrompt(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"one line", "track\n", "gctl> unknown command \"track\"; type help\ngctl> "},
		{"two lines", "track\npoint\n", "gctl> unknown command \"track\"; type help\ngctl> unknown command \"point\"; type help\ngctl> "},
		{"no newline at the end", "track", "gctl> unknown command \"track\"; type help\ngctl> "},
		{"nothing typed", "", "gctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}
