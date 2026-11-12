// Command passes prints the passes of the Apogee world's satellites over its ground stations,
// one a line by AOS: the pass ID, AOS, LOS and maximum elevation. Every pass in the Apogee
// courses comes from this model, so a lesson or a story that needs one takes it from here and
// never invents it (docs/authoring.md, "Calculating the world").
//
//	passes --from <date> [--to <date>] [--sat pipit-2,cirrine-1] [--station YST,PEL]
//
// A date is a day, 2029-07-23, meaning its midnight UTC, or an RFC 3339 instant; --to, not
// included, is a day after --from unless given. Without --sat or --station every satellite or
// station counts.
//
// The model is the one docs/apogee/story.md writes out under "The pass model", copied unchanged
// from the lessons that carry it: model.go from Season 7's, network.go (the twenty stations and
// the eighty satellites' elements) from the capstone's, and types.go, the types they use, from
// Season 7's. Change them only together with the bible and those lessons. The elements hold
// from 1 January 2029, or a satellite's deployment, on: the passes of Seasons 1 to 6 came from
// the same model fitted to earlier epochs, and the bible and the lessons record them.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"
)

const usage = `usage: passes --from <date> [--to <date>] [--sat <ids>] [--station <codes>]
  prints each pass with its AOS in [from, to): pass ID, AOS, LOS and maximum elevation.
  A date is 2029-07-23 (midnight UTC) or an RFC 3339 instant; --to is a day after --from
  unless given. --sat and --station take comma-separated lists; without them, all count.`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run prints the passes args ask for and returns the exit status: 0 for success, 2 for a
// usage mistake.
func run(args []string, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("passes", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fromText := fs.String("from", "", "")
	toText := fs.String("to", "", "")
	satList := fs.String("sat", "", "")
	stationList := fs.String("station", "", "")
	fail := func(format string, a ...any) int {
		fmt.Fprintf(errOut, "passes: "+format+"\n%s\n", append(a, usage)...)
		return 2
	}
	if err := fs.Parse(args); err != nil {
		return fail("%v", err)
	}
	if fs.NArg() > 0 {
		return fail("unexpected argument %q", fs.Arg(0))
	}
	if *fromText == "" {
		return fail("--from is required")
	}
	from, err := parseDate(*fromText)
	if err != nil {
		return fail("--from: %v", err)
	}
	to := from.AddDate(0, 0, 1)
	if *toText != "" {
		if to, err = parseDate(*toText); err != nil {
			return fail("--to: %v", err)
		}
	}
	if !to.After(from) {
		return fail("--to must be after --from")
	}
	n := newNetwork()
	elements, err := pick(n.Elements, *satList, "satellite", func(e Elements) string { return e.SatID })
	if err != nil {
		return fail("%v", err)
	}
	sites, err := pick(n.Sites, *stationList, "station", func(s Site) string { return s.Code })
	if err != nil {
		return fail("%v", err)
	}
	var passes []PassView
	for _, e := range elements {
		if from.Before(e.Epoch) && *satList != "" {
			fmt.Fprintf(errOut, "passes: %s's elements start at %s; it has no passes before that\n", e.SatID, e.Epoch.Format(time.RFC3339))
		}
		for _, s := range sites {
			passes = append(passes, predict(e, s, from, to)...)
		}
	}
	slices.SortStableFunc(passes, func(a, b PassView) int {
		if c := a.AOS.Compare(b.AOS); c != 0 {
			return c
		}
		return strings.Compare(a.ID, b.ID)
	})
	for _, p := range passes {
		fmt.Fprintf(out, "%s %s %s %.1f\n", p.ID, p.AOS.UTC().Format(time.RFC3339), p.LOS.UTC().Format(time.RFC3339), p.MaxElevDeg)
	}
	return 0
}

// parseDate reads a day (midnight UTC) or an RFC 3339 instant.
func parseDate(s string) (time.Time, error) {
	if t, err := time.Parse(time.DateOnly, s); err == nil {
		return t, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("%q is neither a day like 2029-07-23 nor an RFC 3339 instant", s)
	}
	return t.UTC(), nil
}

// pick returns the items list names (comma-separated), in the network's order, or all of them
// for an empty list. A name the network lacks is an error naming every one it has.
func pick[T any](items []T, list, what string, name func(T) string) ([]T, error) {
	if list == "" {
		return items, nil
	}
	want := map[string]bool{}
	for w := range strings.SplitSeq(list, ",") {
		want[strings.TrimSpace(w)] = true
	}
	var picked, all []string
	var out []T
	for _, it := range items {
		all = append(all, name(it))
		if want[name(it)] {
			out = append(out, it)
			picked = append(picked, name(it))
		}
	}
	for w := range want {
		if !slices.Contains(picked, w) {
			return nil, fmt.Errorf("no %s %q; there are %s", what, w, strings.Join(all, ", "))
		}
	}
	return out, nil
}
