package main

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
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
// it wrote before the closing bye with want; TestExit checks the bye.
func check(t *testing.T, input, want string) {
	t.Helper()
	s := newShell()
	s.commands["wind"] = testCommand("read Cape Wren's wind")
	got, err := session(s, input)
	got = strings.TrimSuffix(got, "bye\n")
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
	check(t, "help\n", "gctl> "+
		"exit      leave the shell\n"+
		"help      list every command\n"+
		"passes    list tonight's passes\n"+
		"wind      read Cape Wren's wind\n"+
		"gctl> ")

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
	want := "gctl> " +
		"azimuth   turn the dish to a compass direction\n" +
		"beacon    listen for a satellite's beacon\n" +
		"doppler   show a beacon's frequency shift\n" +
		"exit      leave the shell\n" +
		"help      list every command\n" +
		"passes    list tonight's passes\n" +
		"uplink    send a signed command\n" +
		"wind      read Cape Wren's wind\n" +
		"gctl> "
	if strings.TrimSuffix(got, "bye\n") != want {
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

func TestExit(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"exit", "exit\n", "gctl> bye\n"},
		{"exit with words after it", "exit now\n", "gctl> bye\n"},
		{"lines after exit", "exit\nwind CWR\n", "gctl> bye\n"},
		{"end of input", "wind\n", "gctl> ran with []\ngctl> bye\n"},
		{"nothing typed", "", "gctl> bye\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newShell()
			s.commands["wind"] = testCommand("read Cape Wren's wind")
			got, err := session(s, tc.input)
			if got != tc.want {
				t.Errorf("Run(%q) wrote %q, want %q", tc.input, got, tc.want)
			}
			if err != nil {
				t.Errorf("Run(%q) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestCommandError(t *testing.T) {
	s := newShell()
	s.commands["point"] = command{help: "point the dish", run: func(s *Shell, args []string, out io.Writer) error {
		return errors.New("dish stowed")
	}}
	input := "point CWR\npoint\n"
	got, err := session(s, input)
	want := "gctl> error: dish stowed\ngctl> error: dish stowed\ngctl> bye\n"
	if got != want {
		t.Errorf("Run(%q) wrote %q, want %q", input, got, want)
	}
	if err != nil {
		t.Errorf("Run(%q) = %v, want nil: a command's error is not the session's", input, err)
	}
}

func TestReadError(t *testing.T) {
	problem := errors.New("terminal closed")
	var out strings.Builder
	err := newShell().Run(io.MultiReader(strings.NewReader("track\n"), iotest.ErrReader(problem)), &out)
	if !errors.Is(err, problem) {
		t.Errorf("Run(input failing after one line) = %v, want %v", err, problem)
	}
	want := "gctl> unknown command \"track\"; type help\ngctl> "
	if out.String() != want {
		t.Errorf("Run(input failing after one line) wrote %q, want %q, with no bye", out.String(), want)
	}
}

// tonightLines is what passes prints for the whole night.
const tonightLines = "04:12-04:21 PEL pipit-1 max 61.4\n" +
	"04:31-04:40 YST pipit-2 max 47.2\n" +
	"04:58-05:06 CWR furrow-a max 33.5\n" +
	"05:20-05:26 LAN brine-3 max 22.8\n" +
	"05:47-05:57 YST cirrine-1 max 71.9\n"

func TestPasses(t *testing.T) {
	check(t, "passes\n", "gctl> "+tonightLines+"gctl> ")
	check(t, "passes\npasses\n", "gctl> "+tonightLines+"gctl> "+tonightLines+"gctl> ")
}

func TestPassesAtStation(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"Ystmark", "passes YST\n", "gctl> 04:31-04:40 YST pipit-2 max 47.2\n05:47-05:57 YST cirrine-1 max 71.9\ngctl> "},
		{"Port Ellery", "passes PEL\n", "gctl> 04:12-04:21 PEL pipit-1 max 61.4\ngctl> "},
		{"no passes tonight", "passes SQT\n", "gctl> gctl> "},
		{"unknown code", "passes XYZ\n", "gctl> error: unknown station \"XYZ\"\ngctl> "},
		{"lowercase code", "passes pel\n", "gctl> error: unknown station \"pel\"\ngctl> "},
		{"two codes", "passes PEL YST\n", "gctl> error: passes takes at most one station code\ngctl> "},
		{"the session goes on", "passes XYZ\npasses PEL\n", "gctl> error: unknown station \"XYZ\"\ngctl> 04:12-04:21 PEL pipit-1 max 61.4\ngctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
	s := newShell()
	before := slices.Clone(s.passes)
	session(s, "passes\npasses YST\n")
	if !slices.Equal(s.passes, before) {
		t.Errorf("Run(%q) left s.passes in another order, want it unchanged: sort a copy", "passes\npasses YST\n")
	}
}
