package main

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"
)

// takeOverDay is midnight UTC on Thursday 28 June 2029.
var takeOverDay = time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)

// at returns h:m on take-over day; hours past 24 run into later days.
func at(h, m int) time.Time {
	return takeOverDay.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
}

// station returns n's site with the code.
func station(n *Network, code string) Site {
	for _, s := range n.Sites {
		if s.Code == code {
			return s
		}
	}
	panic("no station " + code)
}

// useModel puts m in passModel's place until the test ends.
func useModel(t *testing.T, m Predictor) {
	before := passModel
	passModel = m
	t.Cleanup(func() { passModel = before })
}

// checkContacts fails the test unless got holds want's contacts, compared by pass ID and dish.
func checkContacts(t *testing.T, call string, got, want []Contact) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s = %d contacts, want %d", call, len(got), len(want))
	}
	for i := range min(len(got), len(want)) {
		if got[i].ID != want[i].ID || got[i].Dish != want[i].Dish {
			t.Errorf("%s[%d] = %s dish %d, want %s dish %d", call, i, got[i].ID, got[i].Dish, want[i].ID, want[i].Dish)
			return
		}
	}
}

// planFor runs plan in a goroutine and waits up to five seconds for it: real predictions, far more time than they take.
func planFor(t *testing.T, call string, plan func() (*Schedule, error)) (*Schedule, error) {
	t.Helper()
	type result struct {
		s   *Schedule
		err error
	}
	done := make(chan result, 1)
	go func() {
		s, err := plan()
		done <- result{s, err}
	}()
	select {
	case r := <-done:
		return r.s, r.err
	case <-time.After(5 * time.Second):
		t.Fatalf("%s did not return", call)
		return nil, nil
	}
}

func TestReplanFrom(t *testing.T) {
	n := newNetwork()
	cfg := Config{Slots: 12, Days: 1, MinElevDeg: 10, Turnaround: 2 * time.Minute}
	day := Window{From: takeOverDay, To: at(24, 0)}
	free := handoverRef(oldPlanner(takeOverDay.AddDate(0, 0, -1)), takeOverDay)
	wind := Stow{Station: "CWR", From: at(9, 47), Reason: "wind 74 km/h"} // no end yet
	calm := wind
	calm.To = at(13, 23)
	v1 := planRef(n, day, free, nil) // 2,016 contacts
	v1.Version = 1
	v60 := planRef(n, day, free, []Stow{wind}) // 1,982 contacts
	v60.Version = 60
	for _, tt := range []struct {
		name  string
		cur   *Schedule
		at    time.Time
		stows []Stow
		want  *Schedule
	}{
		{"version 1 at 09:47, nothing new", v1, at(9, 47), nil, v1},
		{"version 1 at 07:48:33, as Pipit-1 rises over Lantern Atoll", v1, at(7, 48).Add(33 * time.Second), nil, v1},
		{"version 1 at 09:47, Cape Wren stowed with no end yet", v1, at(9, 47), []Stow{wind}, v60},
		{"version 60 at 13:23, the stow ended", v60, at(13, 23), []Stow{calm}, planRef(n, day, free, []Stow{calm})},
	} {
		t.Run(tt.name, func(t *testing.T) {
			call := fmt.Sprintf("replanFrom(version %d, %s)", tt.cur.Version, tt.at.Format(time.TimeOnly))
			before := slices.Clone(tt.cur.Contacts)
			got, err := planFor(t, call, func() (*Schedule, error) {
				return replanFrom(context.Background(), n, tt.cur, tt.at, tt.stows, cfg)
			})
			if err != nil {
				t.Fatalf("%s error = %v, want <nil>", call, err)
			}
			if !slices.EqualFunc(tt.cur.Contacts, before, func(a, b Contact) bool { return a.ID == b.ID && a.Dish == b.Dish }) {
				t.Errorf("%s changed version %d's own contacts, want them as they were: a published schedule never changes", call, tt.cur.Version)
			}
			if got.Window != day {
				t.Errorf("%s window = %s to %s, want version %d's, 28 June 00:00 to 29 June 00:00", call, got.From.Format("2 January 15:04"), got.To.Format("2 January 15:04"), tt.cur.Version)
			}
			checkContacts(t, call, got.Contacts, tt.want.Contacts)
		})
	}
}
