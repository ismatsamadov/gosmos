package main

import (
	"context"
	"errors"
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

func TestPlanStation(t *testing.T) {
	n := newNetwork()
	cfg := Config{Slots: 12, Days: 1, MinElevDeg: 10, Turnaround: 2 * time.Minute}
	day := Window{From: takeOverDay, To: at(24, 0)}
	wind := Stow{Station: "CWR", From: at(9, 47), To: at(13, 23), Reason: "wind 74 km/h"}
	tests := []struct {
		name    string
		call    string // how messages name the call
		code    string
		free    []time.Time // when each dish is free
		stows   []Stow
		dropped int
	}{
		{"Cape Wren on take-over day, its dish free", "planStation(CWR, 28 June)", "CWR", nil, nil, 284},
		{"Cape Wren stowed by the wind from 09:47 to 13:23", "planStation(CWR, 28 June, stowed 09:47 to 13:23)", "CWR", nil, []Stow{wind}, 280},
		{"Saltreach Island, its dish busy until 00:07:00", "planStation(SRH, 28 June, its dish free from 00:07)", "SRH", []time.Time{at(0, 7)}, nil, 144},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := tt.call
			got, dropped, err := planStation(context.Background(), n, station(n, tt.code), day, tt.free, tt.stows, cfg)
			if err != nil {
				t.Fatalf("%s error = %v, want <nil>", call, err)
			}
			want, _ := stationRef(n, station(n, tt.code), day, tt.free, tt.stows)
			checkContacts(t, call, got, want)
			if dropped != tt.dropped {
				t.Errorf("%s dropped %d candidates, want %d", call, dropped, tt.dropped)
			}
		})
	}
	t.Run("a run canceled before it starts", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		got, dropped, err := planStation(ctx, n, station(n, "CWR"), day, nil, nil, cfg)
		if !errors.Is(err, context.Canceled) || len(got) != 0 || dropped != 0 {
			t.Errorf("planStation(CWR, a canceled run) = %d contacts, %d dropped, error %v, want 0, 0 and %v", len(got), dropped, err, context.Canceled)
		}
	})
	t.Run("a week's run canceled when Friday 29 June is first asked for", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var last time.Time // the latest day passModel was asked for
		useModel(t, func(e Elements, s Site, from, to time.Time) []PassView {
			if from.After(last) {
				last = from
			}
			if !from.Before(at(24, 0)) {
				cancel()
			}
			return predict(e, s, from, to)
		})
		week := Window{From: takeOverDay, To: at(7*24, 0)}
		got, _, err := planStation(ctx, n, station(n, "CWR"), week, nil, nil, cfg)
		if last.IsZero() {
			t.Fatalf("planStation(CWR, a week) never asked passModel for a day, want it to: the tests put their own model there")
		}
		if !errors.Is(err, context.Canceled) || len(got) != 0 {
			t.Errorf("planStation(CWR, a week canceled on 29 June) = %d contacts, error %v, want 0 and %v", len(got), err, context.Canceled)
		}
		if last.After(at(4*24, 0)) {
			t.Errorf("planStation(CWR, canceled on 29 June) asked passModel for %s, want no day after 2 July: check ctx at every 1,000th candidate", last.Format("2 January"))
		}
	})
}
