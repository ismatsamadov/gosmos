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
	bookings, err := loadBookings(w.Dir)
	if err != nil {
		return err
	}
	sites := set.sites(n)
	wk := planWeek(n, set.window(), sites, archive)
	booked, unplanned := matchBookings(n, wk, archive, bookings)

	say(w, "rehearsal: %s, %s, %s", set.Week, count(set.Days, "day"), count(len(sites), "station"))
	if err := wk.writePlan(w.Dir); err != nil {
		return err
	}
	say(w, "plan: %d contacts, %d dropped", len(wk.Contacts), wk.Dropped)
	if err := writeUnplanned(w.Dir, unplanned); err != nil {
		return err
	}
	say(w, "bookings: %d, unplanned %d", len(booked), len(unplanned))

	r := &run{w: w, n: n, wk: wk, health: &Health{newest: map[string]TelemetryFrame{}}, outcomes: map[string]Outcome{}}
	r.fly(ctx)
	if ctx.Err() != nil {
		say(w, "rehearsal stopped at %s", w.Now().Format("Mon 15:04:05"))
		return context.Cause(ctx)
	}
	if err := r.writeReports(); err != nil {
		return err
	}
	flown, failed, skipped := r.counts()
	say(w, "week done: %d flown, %d failed, %d skipped", flown, failed, skipped)
	return nil
}

// writeReports saves the reports written after the week: each satellite's newest frame.
func (r *run) writeReports() error {
	return r.health.write(r.w.Dir, r.n)
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
