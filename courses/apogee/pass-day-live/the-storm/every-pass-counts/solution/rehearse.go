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
	masters, err := loadKeys(w.Dir)
	if err != nil {
		return err
	}
	customers, err := loadCustomers(w.Dir)
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

	r, err := newRun(ctx, w, n, wk, booked, masters, customers)
	if err != nil {
		return err
	}
	r.fly(ctx)
	if err := r.logs.save(); err != nil {
		return err
	}
	if ctx.Err() != nil {
		say(w, "rehearsal stopped at %s", w.Now().Format("Mon 15:04:05"))
		return context.Cause(ctx)
	}
	if err := r.writeReports(booked, unplanned); err != nil {
		return err
	}
	flown, failed, skipped := r.counts()
	say(w, "week done: %d flown, %d failed, %d skipped", flown, failed, skipped)
	return nil
}

// newRun prepares the week for flying: every contact's booking, uplinkd, dlinkd, the
// deliveries, the logs and the pager.
func newRun(ctx context.Context, w World, n *Network, wk *Week, bookings []Booking, masters []Master, customers []Customer) (*run, error) {
	r := &run{w: w, n: n, wk: wk, booked: map[string]Booking{}, health: &Health{newest: map[string]TelemetryFrame{}},
		met: newMetrics(), outcomes: map[string]Outcome{}}
	byKey := map[passKey]Booking{}
	for _, b := range bookings {
		byKey[passKey{b.SatID, b.Station, b.AOS.Unix()}] = b
	}
	for _, c := range wk.Contacts {
		if b, ok := byKey[passKey{c.SatID, c.Station, c.AOS.Unix()}]; ok {
			r.booked[c.ID] = b
		}
	}
	r.up = newUplink(r, masters)
	r.down = newDownlink()
	r.deliv = newDeliveries(ctx, r, customers)
	logs, err := newLogs(w, wk.Sites)
	if err != nil {
		return nil, err
	}
	r.logs = logs
	r.pager = newPager(r)
	return r, nil
}

// writeReports saves the reports written after the week: each satellite's newest frame,
// the files still incomplete, the events never delivered, the week's numbers, the alerts
// and the week's report on the bookings given.
func (r *run) writeReports(bookings []Booking, unplanned []Unplanned) error {
	dir := r.w.Dir
	for _, write := range []func() error{
		func() error { return r.health.write(dir, r.n) },
		func() error { return r.down.writeIncomplete(dir) },
		func() error { return r.deliv.writeUndelivered(dir) },
		func() error { return r.met.write(dir) },
		func() error { return r.pager.write(dir) },
		func() error { return r.writeWeek(dir, bookings, unplanned) },
	} {
		if err := write(); err != nil {
			return err
		}
	}
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
