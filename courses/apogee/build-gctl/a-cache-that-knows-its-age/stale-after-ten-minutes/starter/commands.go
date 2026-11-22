package main

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
)

// command is one entry of gctl's command table: its help text and the function that runs it.
type command struct {
	help string
	run  func(s *Shell, args []string, out io.Writer) error
}

// errExit is what the exit command returns to end the session.
var errExit = errors.New("exit")

// commandTable returns every command gctl knows, by name.
func commandTable() map[string]command {
	return map[string]command{
		"exit":     {help: "leave the shell", run: cmdExit},
		"help":     {help: "list every command", run: cmdHelp},
		"passes":   {help: "list tonight's passes", run: cmdPasses},
		"stations": {help: "list the stations and which are stowed", run: cmdStations},
		"stow":     {help: "record a stow: stow CODE [reason...]", run: cmdStow},
	}
}

// cmdHelp lists every command and its help text, sorted by name.
func cmdHelp(s *Shell, args []string, out io.Writer) error {
	for _, name := range slices.Sorted(maps.Keys(s.commands)) {
		fmt.Fprintf(out, "%-8s  %s\n", name, s.commands[name].help)
	}
	return nil
}

// cmdExit ends the session.
func cmdExit(s *Shell, args []string, out io.Writer) error {
	return errExit
}

// cmdPasses lists tonight's passes in order of AOS, at every station or at the one named.
func cmdPasses(s *Shell, args []string, out io.Writer) error {
	if len(args) > 1 {
		return errors.New("passes takes at most one station code")
	}
	if len(args) == 1 && stationIndex(s.stations, args[0]) < 0 {
		return fmt.Errorf("unknown station %q", args[0])
	}
	ps, err := s.cache.passes()
	if err != nil {
		return err
	}
	ps = slices.Clone(ps)
	slices.SortFunc(ps, func(a, b Pass) int { return a.AOS.Compare(b.AOS) })
	for _, p := range ps {
		if len(args) == 0 || p.Station == args[0] {
			fmt.Fprintf(out, "%s-%s %s %s max %.1f\n", p.AOS.Format("15:04"), p.LOS.Format("15:04"), p.Station, p.SatID, p.MaxElevDeg)
		}
	}
	return nil
}

// cmdStations lists every station in the table's order, with its stow if it has one.
func cmdStations(s *Shell, args []string, out io.Writer) error {
	if len(args) > 0 {
		return errors.New("stations takes no arguments")
	}
	for _, st := range s.stations {
		fmt.Fprintf(out, "%s %s%s\n", st.Code, st.Name, stowNote(st))
	}
	return nil
}

// stowNote returns " stowed: <reason>" for a stowed station, " stowed" when no reason was
// given, and "" for a station that is not stowed.
func stowNote(st Station) string {
	switch {
	case !st.Stowed:
		return ""
	case st.StowReason == "":
		return " stowed"
	}
	return " stowed: " + st.StowReason
}

// cmdStow records that a station's dish is stowed, with the reason typed after its code.
func cmdStow(s *Shell, args []string, out io.Writer) error {
	if len(args) == 0 {
		return errors.New("stow needs a station code")
	}
	i := stationIndex(s.stations, args[0])
	if i < 0 {
		return fmt.Errorf("unknown station %q", args[0])
	}
	st := &s.stations[i]
	if st.Stowed {
		return fmt.Errorf("%s is already stowed", st.Code)
	}
	st.Stowed, st.StowReason = true, strings.Join(args[1:], " ")
	fmt.Fprintf(out, "%s%s\n", st.Code, stowNote(*st))
	return nil
}
