package main

import (
	"context"
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

// freeAt returns when dish d+1 of the station is free in free: zero, free now, if free has no entry for it.
func freeAt(free map[string][]time.Time, code string, d int) time.Time {
	if d < len(free[code]) {
		return free[code][d]
	}
	return time.Time{}
}

// clock shows a free time as the tests print it.
func clock(t time.Time) string {
	if t.IsZero() {
		return "free now"
	}
	return t.Format("2 January 15:04:05")
}

// clash returns how many contacts take a dish an old contact holds, until its LOS plus two minutes, and the first of them.
func clash(old, contacts []Contact) (int, string) {
	k, first := 0, ""
	for _, c := range contacts {
		for _, o := range old {
			if o.Station == c.Station && o.Dish == c.Dish && c.AOS.Before(o.LOS.Add(2*time.Minute)) && o.AOS.Before(c.LOS) {
				if k == 0 {
					first = c.ID + " at " + c.AOS.Format(time.TimeOnly) + ", on the dish " + o.ID + " holds until " + o.LOS.Add(2*time.Minute).Format(time.TimeOnly)
				}
				k++
				break
			}
		}
	}
	return k, first
}

func TestHandover(t *testing.T) {
	n := newNetwork()
	old := oldPlanner(takeOverDay.AddDate(0, 0, -1)) // 27 June, 1,653 contacts
	for _, tt := range []struct {
		name string
		at   time.Time
		busy int // dishes still busy at at
	}{
		{"the old planner's 27 June, at 00:00 on take-over day", takeOverDay, 14},
		{"the old planner's 27 June, at 20:40:16 as Furrow-G is first heard", at(-4, 40).Add(16 * time.Second), 11},
	} {
		t.Run(tt.name, func(t *testing.T) {
			call := "handover(27 June, " + tt.at.Format("2 January 15:04:05") + ")"
			got, want := handover(old, tt.at, 2*time.Minute), handoverRef(old, tt.at)
			busy := 0
			for _, s := range n.Sites {
				for d := range s.Dishes {
					if g, w := freeAt(got, s.Code, d), freeAt(want, s.Code, d); !g.Equal(w) {
						t.Fatalf("%s[%q][%d] = %s, want %s", call, s.Code, d, clock(g), clock(w))
					}
					if freeAt(got, s.Code, d).After(tt.at) {
						busy++
					}
				}
			}
			if busy != tt.busy {
				t.Errorf("%s has %d dishes still busy, want %d", call, busy, tt.busy)
			}
		})
	}
	t.Run("Plan from the hand-over at 00:00", func(t *testing.T) {
		cfg := Config{Slots: 12, Days: 1, MinElevDeg: 10, Turnaround: 2 * time.Minute}
		day := Window{From: takeOverDay, To: at(24, 0)}
		free := handover(old, takeOverDay, cfg.Turnaround)
		got, err := planFor(t, "Plan(28 June)", func() (*Schedule, error) { return Plan(context.Background(), n, day, free, nil, cfg) })
		if err != nil {
			t.Fatalf("Plan(28 June, from the hand-over) error = %v, want <nil>", err)
		}
		if k, first := clash(old.Contacts, got.Contacts); k != 0 {
			t.Errorf("Plan(28 June, from the hand-over) gives %d contacts a dish the old planner holds, want 0: %s", k, first)
		}
		checkContacts(t, "Plan(28 June, from the hand-over)", got.Contacts, planRef(n, day, handoverRef(old, takeOverDay), nil).Contacts)
	})
}
