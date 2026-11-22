package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// takeOverDay is midnight UTC on Thursday 28 June 2029, where the tests set a bubble's clock.
var takeOverDay = time.Date(2029, 6, 28, 0, 0, 0, 0, time.UTC)

// clockAt returns a clock that reads day when the bubble starts and moves with its fake time.
func clockAt(day time.Time) func() time.Time {
	shift := day.Sub(time.Now())
	return func() time.Time { return time.Now().Add(shift).UTC() }
}

// at returns h:m on take-over day; hours past 24 run into later days.
func at(h, m int) time.Time {
	return takeOverDay.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
}

// flight is a contact as the station runner flew it.
type flight struct {
	Contact
	at time.Time // when fly was called for it
}

// takeOverUntil runs takeOver in a bubble from 00:00 on take-over day until end, sending news at its times, and returns what flew.
func takeOverUntil(t *testing.T, n *Network, old *Schedule, cfg Config, end time.Time, news []Event) ([]flight, error) {
	var flights []flight
	var err error
	synctest.Test(t, func(t *testing.T) {
		now := clockAt(takeOverDay)
		var sender sync.WaitGroup // sends news at its times
		defer sender.Wait()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		stop := time.AfterFunc(end.Sub(now()), cancel)
		defer stop.Stop()
		events := make(chan Event)
		sender.Go(func() {
			for _, ev := range news {
				select {
				case <-time.After(ev.At.Sub(now())):
				case <-ctx.Done():
					return
				}
				select {
				case events <- ev:
				case <-ctx.Done():
					return
				}
			}
		})
		err = takeOver(ctx, n, old, takeOverDay, events, cfg, now, func(c Contact) {
			flights = append(flights, flight{c, now()})
			if len(flights) > 3000 {
				cancel() // far more than a day holds: end the runner's day
			}
		})
	})
	return flights, err
}

// useModel puts m in passModel's place until the test ends.
func useModel(t *testing.T, m Predictor) {
	before := passModel
	passModel = m
	t.Cleanup(func() { passModel = before })
}

// dayModel returns a pass model that predicts each pair's passes over day once and answers a
// question inside day with those that rise in it, as predict does: the day's replans then
// cost one day of prediction. A question outside day goes to predict.
func dayModel(n *Network, day Window) Predictor {
	type pair struct{ sat, station string }
	passes := map[pair][]PassView{}
	for _, e := range n.Elements {
		for _, s := range n.Sites {
			passes[pair{e.SatID, s.Code}] = predict(e, s, day.From, day.To)
		}
	}
	return func(e Elements, s Site, from, to time.Time) []PassView {
		all, ok := passes[pair{e.SatID, s.Code}]
		if !ok || from.Before(day.From) || to.After(day.To) {
			return predict(e, s, from, to)
		}
		var rise []PassView
		for _, p := range all {
			if !p.AOS.Before(from) && p.AOS.Before(to) {
				rise = append(rise, p)
			}
		}
		return rise
	}
}

// clash returns the first flight on a dish an old contact holds, until its LOS plus two minutes, or "".
func clash(old []Contact, flights []flight) string {
	for _, f := range flights {
		for _, o := range old {
			if o.Station == f.Station && o.Dish == f.Dish && f.AOS.Before(o.LOS.Add(2*time.Minute)) && o.AOS.Before(f.LOS) {
				return fmt.Sprintf("%s at %s on %s's dish %d, which %s holds until %s", f.ID, f.at.Format(time.TimeOnly), f.Station, f.Dish, o.ID, o.LOS.Add(2*time.Minute).Format(time.TimeOnly))
			}
		}
	}
	return ""
}

// checkFlights fails the test unless flights flew want's contacts, each once and at its AOS, and none on a dish the old planner holds.
func checkFlights(t *testing.T, call string, flights []flight, want []Contact, old []Contact) {
	t.Helper()
	seen := map[string]bool{}
	for _, f := range flights {
		if seen[f.ID] {
			t.Errorf("%s flew %s twice, want each contact once", call, f.ID)
			return
		}
		seen[f.ID] = true
		if !f.at.Equal(f.AOS) {
			t.Errorf("%s flew %s at %s, want %s, its AOS", call, f.ID, f.at.Format(time.TimeOnly), f.AOS.Format(time.TimeOnly))
			return
		}
	}
	if c := clash(old, flights); c != "" {
		t.Errorf("%s flew %s, want no dish the old planner holds", call, c)
	}
	if len(flights) != len(want) {
		t.Errorf("%s flew %d contacts, want %d", call, len(flights), len(want))
	}
	for i := range min(len(flights), len(want)) {
		if f, w := flights[i], want[i]; f.ID != w.ID || f.Dish != w.Dish {
			t.Errorf("%s contact %d = %s dish %d at %s, want %s dish %d at %s", call, i+1, f.ID, f.Dish, f.AOS.Format(time.TimeOnly), w.ID, w.Dish, w.AOS.Format(time.TimeOnly))
			return
		}
	}
}

func TestTakeOver(t *testing.T) {
	n := newNetwork()
	old := oldPlanner(takeOverDay.AddDate(0, 0, -1)) // the old planner's 27 June
	free := handoverRef(old, takeOverDay)
	day := Window{From: takeOverDay, To: at(24, 0)}
	useModel(t, dayModel(n, day)) // every replan asks again for passes of the same day
	cfg := Config{Slots: 12, Days: 1, MinElevDeg: 10, Turnaround: 2 * time.Minute, RunLimit: 30 * time.Second}
	t.Run("the cut-over's first 25 minutes, a replan every 10", func(t *testing.T) {
		c := cfg
		c.Replan = 10 * time.Minute
		flights, err := takeOverUntil(t, n, old, c, at(0, 25), nil)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("takeOver(until 00:25) = %v, want %v once its context ends", err, context.Canceled)
		}
		want := planRef(n, Window{From: takeOverDay, To: at(0, 25)}, free, nil).Contacts
		checkFlights(t, "takeOver(until 00:25)", flights, want, old.Contacts)
	})
	t.Run("take-over day, a replan every hour", func(t *testing.T) {
		c := cfg
		c.Replan = time.Hour
		wind := Stow{Station: "CWR", From: at(9, 47), Reason: "wind 74 km/h"}
		calm := wind
		calm.To = at(13, 23)
		news := []Event{{At: wind.From, Stow: &wind}, {At: calm.To, Stow: &calm}}
		flights, err := takeOverUntil(t, n, old, c, at(24, 0).Add(30*time.Second), news)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("takeOver(take-over day) = %v, want %v once its context ends", err, context.Canceled)
		}
		checkFlights(t, "takeOver(take-over day)", flights, planRef(n, day, free, []Stow{calm}).Contacts, old.Contacts)
		furrows, heard := 0, map[string]bool{}
		for i, f := range flights {
			switch {
			case strings.HasPrefix(f.SatID, "furrow-") && f.SatID != "furrow-a" && f.SatID != "furrow-b":
				furrows++
				heard[f.SatID] = true
			case f.ID == "furrow-a/VSK/24529" && (i != 416 || !f.at.Equal(at(5, 29).Add(26*time.Second))):
				t.Errorf("Hollis's pass furrow-a/VSK/24529, booked as bk-7412, was contact %d at %s, want contact 417 at 05:29:26", i+1, f.at.Format(time.TimeOnly))
			case f.ID == "pipit-1/LAN/16797" && !f.at.Equal(at(7, 48).Add(33*time.Second)):
				t.Errorf("pipit-1/LAN/16797, Pipit-1's first contact of the day, flew at %s, want 07:48:33", f.at.Format(time.TimeOnly))
			}
		}
		if len(flights) > 0 && (furrows != 819 || len(heard) != 24) {
			t.Errorf("takeOver(take-over day) flew %d contacts of %d new Furrows, want 819 of all 24", furrows, len(heard))
		}
	})
}
