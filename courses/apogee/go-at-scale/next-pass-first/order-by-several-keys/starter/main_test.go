package main

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
	"time"
)

// pass returns a predicted pass from its ID, its AOS and LOS ("2029-04-08 10:29:34", UTC) and its max elevation.
func pass(id, aos, los string, elevDeg float64) PassView {
	sat, rest, _ := strings.Cut(id, "/")
	station, _, _ := strings.Cut(rest, "/")
	at := func(s string) time.Time {
		t, _ := time.Parse(time.DateTime, s)
		return t
	}
	return PassView{ID: id, Pass: Pass{SatID: sat, Station: station, AOS: at(aos), LOS: at(los), MaxElevDeg: elevDeg}}
}

// portEllery is Port Ellery's passes on Sunday 8 April from 09:53 to 11:35, by AOS and then ID.
var portEllery = []PassView{
	pass("furrow-b/PEL/23297", "2029-04-08 09:53:24", "2029-04-08 10:01:12", 19.4),
	pass("brine-3/PEL/16701", "2029-04-08 10:29:34", "2029-04-08 10:38:49", 67.7),
	pass("pipit-3/PEL/11171", "2029-04-08 10:29:34", "2029-04-08 10:35:05", 10.1),
	pass("furrow-a/PEL/23300", "2029-04-08 10:39:14", "2029-04-08 10:48:28", 86.6),
	pass("cirrine-1/PEL/18660", "2029-04-08 10:45:40", "2029-04-08 10:54:50", 25.7),
	pass("furrow-b/PEL/23298", "2029-04-08 11:26:48", "2029-04-08 11:35:28", 33.1),
}

// Two passes over Lantern Atoll that rise in the same second on 11 June.
var (
	brine34   = pass("brine-34/LAN/00403", "2029-06-11 19:09:47", "2029-06-11 19:14:51", 13.6)
	cirrine10 = pass("cirrine-10/LAN/00827", "2029-06-11 19:09:47", "2029-06-11 19:14:09", 11.2)
)

// sign names the sign of a comparison's result as the messages print it.
func sign(n int) string {
	switch {
	case n < 0:
		return "below 0"
	case n > 0:
		return "above 0"
	}
	return "0"
}

func TestByAOS(t *testing.T) {
	brine3, pipit3, furrowA := portEllery[1], portEllery[2], portEllery[3]
	tests := []struct {
		name string
		a, b PassView
		want int // the sign byAOS gives
	}{
		{"Furrow-A rises after Brine-3", furrowA, brine3, +1},
		{"Pipit-3 rises before Furrow-A, whatever their IDs", pipit3, furrowA, -1},
		{"Brine-3, then Pipit-3: both rise at 10:29:34", brine3, pipit3, -1},
		{"Pipit-3, then Brine-3: both rise at 10:29:34", pipit3, brine3, +1},
		{"Brine-34 and Cirrine-10, both at 19:09:47 on 11 June", brine34, cirrine10, -1},
		{"the same pass twice", brine3, brine3, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byAOS(tt.a, tt.b); sign(got) != sign(tt.want) {
				t.Errorf("byAOS(%s, %s) = %d, want %s", tt.a.ID, tt.b.ID, got, sign(tt.want))
			}
		})
	}
}

func TestSortByAOS(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	deck := slices.Clone(portEllery)
	for i := range 3 {
		r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
		passes := slices.Clone(deck)
		slices.SortFunc(passes, byAOS)
		for j, p := range passes {
			if p.ID != portEllery[j].ID {
				t.Fatalf("Sunday's passes in shuffle %d, sorted by byAOS: pass %d = %s, want %s", i+1, j, p.ID, portEllery[j].ID)
			}
		}
	}
}
