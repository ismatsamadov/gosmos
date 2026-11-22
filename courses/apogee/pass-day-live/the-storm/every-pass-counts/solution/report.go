package main

import (
	"bytes"
	"cmp"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

// lostBooking is a booked pass that wasn't flown, and how.
type lostBooking struct {
	Booking
	how string // "skipped stowed", "failed line_lost", "unplanned no such pass"
}

// target is a station's availability target in basis points of booked passes flown:
// 99.0% at the atolls, whose palms and rain fade cost about a pass in a hundred, and 99.5%
// everywhere else.
func target(code string) int {
	if code == "LAN" || code == "CAL" {
		return 9900
	}
	return 9950
}

// writeWeek saves reports/week.txt in dir: a row per rehearsed station and one for all of
// them, a line per lost booking and the verdict.
func (r *run) writeWeek(dir string, bookings []Booking, unplanned []Unplanned) error {
	why := map[string]string{}
	for _, u := range unplanned {
		why[u.ID] = "unplanned " + u.Reason
	}
	passOf := map[string]string{} // booking ID to the pass ID it matched
	for id, b := range r.booked {
		passOf[b.ID] = id
	}
	var lost []lostBooking
	booked, lostAt := map[string]int{}, map[string]int{}
	for _, b := range bookings {
		if !slices.ContainsFunc(r.wk.Sites, func(s Site) bool { return s.Code == b.Station }) {
			continue
		}
		booked[b.Station]++
		how, ok := why[b.ID]
		if !ok {
			if o := r.outcomeOf(passOf[b.ID]); o.Kind != "flown" {
				how = o.String()
			}
		}
		if how != "" {
			lostAt[b.Station]++
			lost = append(lost, lostBooking{b, how})
		}
	}
	y, wk := r.wk.Window.From.ISOWeek()
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%-8s %7s %5s %6s %7s %5s\n", fmt.Sprintf("%d-W%02d", y, wk), "booked", "lost", "avail", "budget", "left")
	var all [3]int
	for _, s := range r.wk.Sites {
		bk, l := booked[s.Code], lostAt[s.Code]
		budget := bk * (10000 - target(s.Code)) / 10000
		row(&buf, s.Code, bk, l, budget)
		all[0], all[1], all[2] = all[0]+bk, all[1]+l, all[2]+budget
	}
	row(&buf, "all", all[0], all[1], all[2])
	slices.SortFunc(lost, func(a, b lostBooking) int { return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID)) })
	for _, l := range lost {
		fmt.Fprintf(&buf, "%s %s %s/%s %s\n", l.ID, l.AOS.UTC().Format("Mon 15:04:05"), l.SatID, l.Station, l.how)
	}
	switch len(lost) {
	case 0:
		buf.WriteString("GO: every booked pass flown\n")
	case 1:
		buf.WriteString("NO-GO: 1 booked pass lost\n")
	default:
		fmt.Fprintf(&buf, "NO-GO: %d booked passes lost\n", len(lost))
	}
	return saveWhole(filepath.Join(dir, "reports", "week.txt"), buf.Bytes())
}

// row writes one row of the week's report; availability is "-" with nothing booked.
func row(buf *bytes.Buffer, code string, booked, lost, budget int) {
	if booked == 0 {
		fmt.Fprintf(buf, "%-8s %7d %5d %6s %7d %5d\n", code, booked, lost, "-", budget, budget-lost)
		return
	}
	avail := 100 * float64(booked-lost) / float64(booked)
	fmt.Fprintf(buf, "%-8s %7d %5d %6.2f %7d %5d\n", code, booked, lost, avail, budget, budget-lost)
}

// outcomeOf returns how a contact ended.
func (r *run) outcomeOf(passID string) Outcome {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.outcomes[passID]
}
