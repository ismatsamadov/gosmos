package main

import (
	"context"
	"slices"
	"strings"
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

// stowsShown describes stows as the tests print them.
func stowsShown(stows []Stow) string {
	var s []string
	for _, st := range stows {
		if st.To.IsZero() {
			s = append(s, st.Station+" from "+st.From.Format("15:04")+", no end")
		} else {
			s = append(s, st.Station+" from "+st.From.Format("15:04")+" to "+st.To.Format("15:04"))
		}
	}
	return "[" + strings.Join(s, "; ") + "]"
}

func TestReplannerStows(t *testing.T) {
	n := newNetwork()
	cfg := Config{Slots: 12, Days: 1, MinElevDeg: 10, Turnaround: 2 * time.Minute, Replan: 10 * time.Minute, RunLimit: 30 * time.Second}
	day := Window{From: takeOverDay, To: at(24, 0)}
	free := handoverRef(oldPlanner(takeOverDay.AddDate(0, 0, -1)), takeOverDay)
	v1 := planRef(n, day, free, nil) // 2,016 contacts
	v1.Version = 1
	pub := &Publisher{}
	pub.Offer(v1)
	r := newReplanner(n, pub, cfg, func() time.Time { return takeOverDay })
	wind := Stow{Station: "CWR", From: at(9, 47), Reason: "wind 74 km/h"}
	calm := wind
	calm.To = at(13, 23)
	lost := []string{"pipit-2/CWR/22443", "cirrine-1/CWR/19868", "furrow-b/CWR/24530", "cirrine-1/CWR/19869", "furrow-g/CWR/00011", "furrow-d/CWR/00011"}
	for _, step := range []struct {
		name    string
		ev      Event
		version uint64 // the current schedule's, after the event's run
		stows   []Stow // the stows heard of, after it
		cwr     int    // Cape Wren's contacts on the day, after it
	}{
		{"a tick at 09:00, with no news", Event{At: at(9, 0)}, 1, nil, 76},
		{"Cape Wren's wind at 09:47, with no end yet", Event{At: at(9, 47), Stow: &wind}, 3, []Stow{wind}, 42},
		{"the wind's end at 13:23", Event{At: at(13, 23), Stow: &calm}, 4, []Stow{calm}, 71},
	} {
		ok := t.Run(step.name, func(t *testing.T) {
			before, kept := r.stows, slices.Clone(r.stows) // what earlier runs read
			events := make(chan Event, 1)
			events <- step.ev
			close(events)
			done := make(chan error, 1)
			go func() { done <- r.Run(context.Background(), events) }()
			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("Run = %v once events closed, want <nil>", err)
				}
			case <-time.After(5 * time.Second): // a real day of planning: far more time than it takes
				t.Fatalf("Run did not return once events closed")
			}
			if !slices.Equal(r.stows, step.stows) {
				t.Errorf("stows heard of = %s, want %s", stowsShown(r.stows), stowsShown(step.stows))
			}
			if !slices.Equal(before, kept) {
				t.Errorf("the stows an earlier run read became %s, want %s still: give r.stows a new slice", stowsShown(before), stowsShown(kept))
			}
			cur := pub.Current()
			if cur.Version != step.version {
				t.Fatalf("current schedule = run %d's, want run %d's", cur.Version, step.version)
			}
			cwr, flies := 0, ""
			for _, c := range cur.Contacts {
				if c.Station == "CWR" {
					cwr++
				}
				if slices.Contains(lost, c.ID) && len(step.stows) > 0 && flies == "" {
					flies = c.ID + " at " + c.AOS.Format(time.TimeOnly)
				}
			}
			if cwr != step.cwr {
				t.Errorf("run %d gives Cape Wren %d contacts, want %d", cur.Version, cwr, step.cwr)
			}
			if flies != "" {
				t.Errorf("run %d flies %s, want it lost to the wind", cur.Version, flies)
			}
			checkContacts(t, "the current schedule", cur.Contacts, planRef(n, day, free, step.stows).Contacts)
		})
		if !ok {
			break // the next step builds on this one
		}
	}
}
