package main

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Unplanned is a booked pass the plan can't fly, and why.
type Unplanned struct {
	Booking
	Reason string // "no such pass"
}

// passKey is what a booking knows of its pass: the satellite, the station and the AOS, to
// the second. A booking carries no pass ID.
type passKey struct {
	SatID, Station string
	AOS            int64 // Unix seconds
}

// loadBookings reads bookings.json in dir strictly: bookingd's export of the week's
// bookings. A missing file means no bookings.
func loadBookings(dir string) ([]Booking, error) {
	f, err := os.Open(filepath.Join(dir, "bookings.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var bookings []Booking
	if err := decodeStrict(f, &bookings, "bookings"); err != nil {
		return nil, fmt.Errorf("bookings.json: %w", err)
	}
	return bookings, nil
}

// matchBookings sets the bookings against the plan. It returns the bookings the rehearsal
// answers for and, by AOS then booking ID, those no contact matches, each with its reason.
func matchBookings(n *Network, wk *Week, archive []PassView, bookings []Booking) ([]Booking, []Unplanned) {
	planned := map[passKey]bool{}
	for _, c := range wk.Contacts {
		planned[passKey{c.SatID, c.Station, c.AOS.Unix()}] = true
	}
	var archived map[passKey]PassView // made when a booking first needs it
	var kept []Booking
	var unplanned []Unplanned
	for _, b := range bookings {
		if !wk.answersFor(n, b) {
			continue
		}
		kept = append(kept, b)
		if key := (passKey{b.SatID, b.Station, b.AOS.Unix()}); !planned[key] {
			if archived == nil {
				archived = map[passKey]PassView{}
				for _, p := range archive {
					archived[passKey{p.SatID, p.Station, p.AOS.Unix()}] = p
				}
			}
			p, ok := archived[key]
			unplanned = append(unplanned, Unplanned{b, whyUnplanned(n, b, p, ok)})
		}
	}
	slices.SortFunc(unplanned, func(a, b Unplanned) int {
		return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
	})
	return kept, unplanned
}

// answersFor reports whether the rehearsal answers for b: its AOS is in the window, and
// its station is one the rehearsal flies or one the network doesn't know.
func (wk *Week) answersFor(n *Network, b Booking) bool {
	if b.AOS.Before(wk.Window.From) || !b.AOS.Before(wk.Window.To) {
		return false
	}
	isStation := func(s Site) bool { return s.Code == b.Station }
	return slices.ContainsFunc(wk.Sites, isStation) || !slices.ContainsFunc(n.Sites, isStation)
}

// whyUnplanned gives the first reason that fits a booking no contact matches; p is the
// archived pass with the booking's satellite, station and AOS, if archived says there is one.
func whyUnplanned(n *Network, b Booking, p PassView, archived bool) string {
	sat := satellite(n, b.SatID)
	site := slices.IndexFunc(n.Sites, func(s Site) bool { return s.Code == b.Station })
	switch {
	case sat.ID == "":
		return "unknown satellite"
	case site < 0:
		return "unknown station"
	case !archived:
		return "no such pass"
	case !canServe(n.Sites[site], sat):
		return "not servable"
	case p.MaxElevDeg < minElevDeg:
		return "too low"
	}
	return "no dish free"
}

// writeUnplanned saves reports/unplanned.txt in dir: one line per unplanned booking, as
// "bk-9902 Thu 04:58:30 furrow-a/CWR no such pass".
func writeUnplanned(dir string, unplanned []Unplanned) error {
	var buf bytes.Buffer
	for _, u := range unplanned {
		fmt.Fprintf(&buf, "%s %s %s/%s %s\n", u.ID, u.AOS.UTC().Format("Mon 15:04:05"), u.SatID, u.Station, u.Reason)
	}
	return saveWhole(filepath.Join(dir, "reports", "unplanned.txt"), buf.Bytes())
}
