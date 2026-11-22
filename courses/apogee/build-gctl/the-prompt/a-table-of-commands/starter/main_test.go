package main

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

// session runs s on input and returns what it wrote and what Run returned.
func session(s *Shell, input string) (string, error) {
	var out strings.Builder
	err := s.Run(strings.NewReader(input), &out)
	return out.String(), err
}

// testCommand returns a test-only command that writes its arguments.
func testCommand(help string) command {
	return command{help: help, run: func(s *Shell, args []string, out io.Writer) error {
		fmt.Fprintf(out, "ran with %q\n", args)
		return nil
	}}
}

// check runs input on a new shell that also knows a test-only wind command, and compares what
// it wrote with want.
func check(t *testing.T, input, want string) {
	t.Helper()
	s := newShell()
	s.commands["wind"] = testCommand("read Cape Wren's wind")
	got, err := session(s, input)
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

func TestWords(t *testing.T) {
	unknown := "gctl> unknown command \"track\"; type help\ngctl> "
	for _, tc := range []struct{ name, input, want string }{
		{"two spaces", "track  pipit-1\n", unknown},
		{"tabs and spaces", "\t track \tpipit-1   PEL \n", unknown},
		{"blank line", "\n", "gctl> gctl> "},
		{"spaces and tabs only", "  \t \n", "gctl> gctl> "},
		{"blank line between", "track\n\npoint\n", "gctl> unknown command \"track\"; type help\ngctl> gctl> unknown command \"point\"; type help\ngctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}

func TestHelp(t *testing.T) {
	check(t, "help\n", "gctl> help      list every command\nwind      read Cape Wren's wind\ngctl> ")

	s := newShell()
	for _, c := range []struct{ name, help string }{
		{"wind", "read Cape Wren's wind"},
		{"uplink", "send a signed command"},
		{"doppler", "show a beacon's frequency shift"},
		{"beacon", "listen for a satellite's beacon"},
		{"azimuth", "turn the dish to a compass direction"},
	} {
		s.commands[c.name] = testCommand(c.help)
	}
	got, _ := session(s, "help please\n")
	want := "gctl> azimuth   turn the dish to a compass direction\n" +
		"beacon    listen for a satellite's beacon\n" +
		"doppler   show a beacon's frequency shift\n" +
		"help      list every command\n" +
		"uplink    send a signed command\n" +
		"wind      read Cape Wren's wind\n" +
		"gctl> "
	if got != want {
		t.Errorf("Run(%q) with five test commands wrote %q, want %q", "help please\n", got, want)
	}
}

func TestCommands(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"arguments", "wind CWR  gusts\n", "gctl> ran with [\"CWR\" \"gusts\"]\ngctl> "},
		{"no arguments", "wind\n", "gctl> ran with []\ngctl> "},
		{"unknown command", "track CWR\n", "gctl> unknown command \"track\"; type help\ngctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}
