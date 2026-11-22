package main

import (
	"context"
	"fmt"
)

// rehearse runs the week rehearsal.json describes against w and returns when every pass,
// file, delivery and report is done; when ctx ends first it stops and returns
// context.Cause(ctx).
func rehearse(ctx context.Context, w World) error {
	n := newNetwork()
	set, err := loadSettings(w.Dir, n)
	if err != nil {
		return err
	}
	archive, err := loadArchive(w.Dir)
	if err != nil {
		return err
	}
	sites := set.sites(n)
	wk := planWeek(n, set.window(), sites, archive)

	say(w, "rehearsal: %s, %s, %s", set.Week, count(set.Days, "day"), count(len(sites), "station"))
	if err := wk.writePlan(w.Dir); err != nil {
		return err
	}
	say(w, "plan: %d contacts, %d dropped", len(wk.Contacts), wk.Dropped)
	return nil
}

// say writes one line on the console, in one Write.
func say(w World, format string, args ...any) {
	fmt.Fprintf(w.Out, format+"\n", args...)
}

// count writes n with its word, in the singular for one: "1 day", "7 days".
func count(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return fmt.Sprintf("%d %ss", n, word)
}
