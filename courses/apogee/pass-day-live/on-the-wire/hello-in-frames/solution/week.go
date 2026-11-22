package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
)

const (
	minElevDeg = 10              // a pass must climb this high to be a candidate
	turnaround = 2 * time.Minute // a dish is free this long after its last contact's LOS
)

// Week is the planned week: what the rehearsal flies.
type Week struct {
	Window   Window
	Sites    []Site    // the stations it flies, in table order
	Contacts []Contact // every contact, by AOS then pass ID
	Dropped  int       // candidates no dish could take
}

// planWeek plans the window at the sites from the archive by passplan 2's rule: each station
// on its own, its candidates in AOS order, each taking the lowest-numbered dish free at its
// AOS.
func planWeek(n *Network, win Window, sites []Site, archive []PassView) *Week {
	wk := &Week{Window: win, Sites: sites}
	for _, site := range sites {
		contacts, dropped := assign(site, slices.Values(candidates(n, site, archive)), nil, turnaround)
		wk.Contacts = append(wk.Contacts, contacts...)
		wk.Dropped += dropped
	}
	slices.SortFunc(wk.Contacts, func(a, b Contact) int { return byAOS(a.PassView, b.PassView) })
	return wk
}

// candidates returns the archive's passes over site that it can serve and that climb to at
// least 10 degrees, by AOS then pass ID.
func candidates(n *Network, site Site, archive []PassView) []PassView {
	var out []PassView
	for _, p := range archive {
		if p.Station == site.Code && p.MaxElevDeg >= minElevDeg && canServe(site, satellite(n, p.SatID)) {
			out = append(out, p)
		}
	}
	slices.SortFunc(out, byAOS)
	return out
}

// satellite returns the network's satellite with the given ID, or one with no bands.
func satellite(n *Network, id string) Satellite {
	for _, s := range n.Satellites {
		if s.ID == id {
			return s
		}
	}
	return Satellite{}
}

// writePlan saves plan/<day>.txt in dir for every day of the window: the contacts whose
// AOS falls on that UTC day, as the dish loaders read them.
func (wk *Week) writePlan(dir string) error {
	for day := wk.Window.From; day.Before(wk.Window.To); day = day.Add(24 * time.Hour) {
		var buf bytes.Buffer
		for _, c := range wk.Contacts {
			if !c.AOS.Before(day) && c.AOS.Before(day.Add(24*time.Hour)) {
				fmt.Fprintf(&buf, "%s %s %d %s\n", c.AOS.Format(time.TimeOnly), c.LOS.Format(time.TimeOnly), c.Dish, c.ID)
			}
		}
		if err := saveWhole(filepath.Join(dir, "plan", day.Format(time.DateOnly)+".txt"), buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// saveWhole writes data to a private temporary file beside path, then renames it over
// path, so a reader sees all of the old file or all of the new one. It makes path's
// folder first.
func saveWhole(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
